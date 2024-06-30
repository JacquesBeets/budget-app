package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RenameFile() error {
	// Read the contents of the unzipped directory
	files, err := os.ReadDir("./unzipped")
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
		return err
	}

	// Check if there's at least one file in the directory
	if len(files) > 0 {
		// Get the first file
		file := files[0]

		// Check if the file is not a directory
		if !file.IsDir() {
			// Rename the file
			oldPath := "./unzipped/" + file.Name()
			newPath := "./unzipped/fnb.ofx"
			err = os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("Failed to rename file: %v", err)
				return err
			}
		}
	}

	return nil
}

func CleanUpDownloads() error {
	err := os.RemoveAll("./downloads")
	if err != nil {
		log.Printf("could not remove downloads directory: %v", err)
		return err
	}

	err = os.RemoveAll("./unzipped")
	if err != nil {
		log.Printf("could not remove unzipped directory: %v", err)
		return err
	}

	return nil
}

func StartDownload() {
	err := playwright.Install()
	if err != nil {
		assertErrorToNilf("could not install playwright: %w", err)
		return
	}
	DownloadFnb()
	DownloadNed()
}

func DownloadNed() error {
	usern := os.Getenv("USERN")
	pass := os.Getenv("PASSWORD")
	website := os.Getenv("WEBSITE")
	waitForLogin := os.Getenv("WEBSITE_LOGIN_WAIT")
	waitForLogout := os.Getenv("WEBSITE_LOGOUT_WAIT")

	// err := playwright.Install()
	// assertErrorToNilf("could not install playwright: %w", err)

	// Launch Playwright
	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)

	// Launch Browser with UI
	// browser, err := pw.Chromium.Launch()
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	// browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
	// 	Headless: playwright.Bool(true),
	// })
	assertErrorToNilf("could not launch Chromium: %w", err)

	defer func() {
		browser.Close()
		// pw.Stop()
	}()

	// Create New Page
	page, err := browser.NewPage()
	assertErrorToNilf("could not create page: %w", err)

	// Goto Website
	_, err = page.Goto(website)
	assertErrorToNilf("could not goto: %w", err)

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not select Use Nedbank ID to log in: %v", page.Locator(`[aria-label="Use Nedbank ID to log in"]`).Click())

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Fill in Username
	assertErrorToNilf("could not type: %v", page.Locator("input#username").Fill(usern))

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Fill in Password
	assertErrorToNilf("could not type: %v", page.Locator("input#password").Fill(pass))

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click Login
	assertErrorToNilf("could not press: %v", page.Locator("#log_in").Click())

	//WaitForLogin to complete
	frame := page.MainFrame()
	_ = frame.WaitForURL(waitForLogin)

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	assertErrorToNilf("could not select statement-position: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-landing/div[1]/div/div[2]/div/div/div/div[1]/a`).Click())

	time.Sleep(3 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not select statement enquiry tab: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-statement-documents-global/div/section/section[1]/app-toggle-tab-group/div/div[2]/label`).Click())

	time.Sleep(3 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not open Format dropdown: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-statement-documents-global/div/section/section[2]/div/app-statements-enquiry/form/div/app-enquiry-dropdown/div/div`).Click())

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not select OFX option: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-statement-documents-global/div/section/section[2]/div/app-statements-enquiry/form/div/app-enquiry-dropdown/div/ul/li[3]`).Click())

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	//Download
	download, err := page.ExpectDownload(func() error {
		return page.Locator(`#download`).Click()
	})
	assertErrorToNilf("could not download file:  %w", err)

	// Save download to file
	err = download.SaveAs("./downloads/ned.ofx") // Save to current directory
	assertErrorToNilf("could not save download to file: %w", err)

	// Logout
	assertErrorToNilf("could not logout: %v", page.Locator("header .shiftHeader li.logout a").Click())

	// WaitForLogout to complete
	_ = frame.WaitForURL(waitForLogout)

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Parse OFX
	err = ParseOFX("./downloads/ned.ofx", "NEDBANK")
	assertErrorToNilf("could not parse OFX: %w", err)

	// Clean Up
	err = CleanUpDownloads()
	assertErrorToNilf("could not clean up: %w", err)

	log.Printf("Completed NED Download")
	return nil
}

func DownloadFnb() error {

	usern := os.Getenv("FNB_USERN")
	pass := os.Getenv("FNB_PASSWORD")
	website := os.Getenv("FNB_WEBSITE")
	waitForLogin := os.Getenv("FNB_WAIT_FOR_LOGIN")
	waitForLogout := os.Getenv("FNB_WAIT_FOR_LOGOUT")

	// Launch Playwright
	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)

	// Launch Browser
	// browser, err := pw.Chromium.Launch()

	// Luanch Browser with UI
	// browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
	// 	Headless: playwright.Bool(false),
	// })
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})

	defer func() {
		browser.Close()
		// pw.Stop()
	}()
	assertErrorToNilf("could not launch Chromium: %w", err)

	// Create New Page
	page, err := browser.NewPage()
	assertErrorToNilf("could not create page: %w", err)

	// Create New Page with Video Recording
	// page, err := browser.NewPage(playwright.BrowserNewPageOptions{
	// 	RecordVideo: &playwright.RecordVideo{
	// 		Dir: "videos/",
	// 	},
	// })
	// assertErrorToNilf("could not create page: %w", err)
	// _, err = page.Video().Path()
	// assertErrorToNilf("failed to get video path: %v", err)

	// Goto Website
	_, err = page.Goto(website)
	assertErrorToNilf("could not goto: %w", err)

	// Fill in Username
	assertErrorToNilf("could not type: %v", page.Locator("input#user").Fill(usern))

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Fill in Password
	assertErrorToNilf("could not type: %v", page.Locator("input#pass").Fill(pass))

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click Login
	assertErrorToNilf("could not press: %v", page.Locator("#OBSubmit").Press("Enter"))

	//WaitForLogin to complete
	frame := page.MainFrame()
	_ = frame.WaitForURL(waitForLogin)

	// time.Sleep(5 * time.Second) // Wait for 5 seconds

	// Click on Accounts
	assertErrorToNilf("could not Click on Accounts: %v", page.Locator("#shortCutLinks > span:nth-child(1)").Click())

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click on Balance
	assertErrorToNilf("could not Click on Balance: %v", page.Locator("#tabelRow_6 .group3 .col4 a").Click())

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click More
	assertErrorToNilf("could not Click More: %v", page.Locator("#footerButtonsContainer > div:nth-child(1) a").Click())

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click on Download
	assertErrorToNilf("could not Click on Download Button: %v", page.Locator("#tableActionButtons .downloadButton").Click())

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Open Dropdown
	assertErrorToNilf("could not open dropdown: %v", page.Locator("#downloadFormat_dropId").Click())

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Select OFX
	assertErrorToNilf("could not select OFX: %v", page.Locator(`[data-value="ofx"]`).Click())
	// assertErrorToNilf("could not select OFX: %v", page.Locator("ul.dropdown-content li:last-child").Click())
	// assertErrorToNilf("could not select OFX: %v", page.Locator("//*[@id="downloadFormat_parent"]/div[2]/div[3]/ul/li[6]").Click())  // X-PATH

	// time.Sleep(3 * time.Second) // Wait for 3 seconds

	//Download
	download, err := page.ExpectDownload(func() error {
		return page.Locator("#eziPannelButtonsWrapper #mainDownloadBtn").Click()
	})
	assertErrorToNilf("could not download file:  %w", err)

	// Save download to file
	err = download.SaveAs("./downloads/fnb_ofx.zip") // Save to current directory
	assertErrorToNilf("could not save download to file: %w", err)

	time.Sleep(5 * time.Second) // Wait for 5 seconds

	// Logout
	assertErrorToNilf("could not logout: %v", page.Locator("#headerButton_").Click())

	// WaitForLogout to complete
	_ = frame.WaitForURL(waitForLogout)

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Unzip
	err = Unzip("./downloads/fnb_ofx.zip", "./unzipped")
	assertErrorToNilf("could not unzip: %w", err)

	// Rename File
	err = RenameFile()
	assertErrorToNilf("could not rename file: %w", err)

	// Parse OFX
	err = ParseOFX("./unzipped/fnb.ofx", "FNB")
	assertErrorToNilf("could not parse OFX: %w", err)

	// Clean Up
	// err = CleanUpDownloads()
	// assertErrorToNilf("could not clean up: %w", err)

	log.Println("Completed FNB Download")
	return nil
}

func StartNewTrace() {
	/// Launch Playwright
	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)
	// Launch Browser with UI
	// browser, err := pw.Chromium.Launch()
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	assertErrorToNilf("could not launch Chromium: %w", err)

	defer func() {
		browser.Close()
		pw.Stop()
	}()

	// Create New Page
	page, err := browser.NewPage()
	assertErrorToNilf("could not create page: %w", err)

	if _, err = page.Goto("http://localhost:9090/"); err != nil {
		log.Fatalf("could not goto: %v\n", err)
	}

	// TODO: Think about possibly injecting a script to listen for clicks and then send the object to the server
	// We could possibly also post to the server and then the server can handle the rest
	// npx playwright codegen
	handle, err := page.EvaluateHandle(`
	function ListenClick() {
		return new Promise(function(resolve, reject) {
			// this works but we need to populate an object of clicks the send it to the server
			// perhaps we can inject a control to stop listening for clicks and then send the object
	
			let clicks = [];
	
			// Create a button element
			var button = document.createElement('button');
			button.innerText = 'Stop Recording';
	
			// position the button and style it
			button.style.position = 'fixed';
			button.style.top = '10px';
			button.style.right = '10px';
			button.style.zIndex = '9999';
			button.style.backgroundColor = 'red';
			button.style.color = 'white';
	
			// Append the button to the body
			document.body.appendChild(button);
	
			// Listen for click events on the button
			button.addEventListener('click', function() {
				// Resolve the promise with a custom value
				resolve(clicks);
			});
	
			function getXPath(element) {
				if (element.id !== '') {
				  return 'id("' + element.id + '")';
				}
				if (element === document.body) {
				  return element.tagName;
				}
			  
				var ix = 0;
				var siblings = element.parentNode.childNodes;
				for (var i = 0; i < siblings.length; i++) {
				  var sibling = siblings[i];
				  if (sibling === element) {
					return getXPath(element.parentNode) + '/' + element.tagName + '[' + (ix + 1) + ']';
				  }
				  if (sibling.nodeType === 1 && sibling.tagName === element.tagName) {
					ix++;
				  }
				}
			  }
		
		
			// Listen for click events on the entire document
			document.addEventListener('click', function(event) {
			  // Get the XPath of the clicked element
			  var xpath = getXPath(event.target);
			  
			  // Log the XPath
			  console.log('Clicked element XPath: ' + xpath);
	
			  clicks.push(xpath);
			});
		  });
	}
	`, struct{}{})
	if err != nil {
		log.Fatalf("could not acquire JSHandle: %v\n", err)
	}

	jsonValue, err := handle.JSONValue()
	if err != nil {
		log.Fatalf("could not get JSON value: %v\n", err)
	}
	fmt.Println("handle", jsonValue)

}
