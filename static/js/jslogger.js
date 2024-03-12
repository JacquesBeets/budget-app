function ListenClick() {
    return new Promise(function(resolve, reject) {
        // TODO: Create a key logger that logs the keys pressed


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
