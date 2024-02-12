package models


type TransactionType struct {
    ID                int       `json:"id"`
    Title             string    `json:"title"`
}
    

func (t *TransactionType) Save() error {
    return nil
}

func (t *TransactionType) Update() error {
    return nil
}

func (t *TransactionType) Delete() error {
    return nil
}