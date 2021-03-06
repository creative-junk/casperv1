//Handler methods for our endpoints
package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Controller struct {
	repository Repository
}

//NewBusiness POST
func (c *Controller) NewBusiness(w http.ResponseWriter, r *http.Request) {
	var business Business
	//read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Fatalln("Error Setting Up Business", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Setting Up Business", err)
	}

	if err := json.Unmarshal(body, &business); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	success := c.repository.setupBusiness(business) //Add Business to Database

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}

// UpdateInvoice PUT
func (c *Controller) ModifyBusiness(w http.ResponseWriter, r *http.Request) {
	var business Business
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Business", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Business", err)
	}

	if err := json.Unmarshal(body, &business); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.updateBusiness(business) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//Invoices /invoices GET
func (c *Controller) ListInvoice(w http.ResponseWriter, r *http.Request) {
	invoices := c.repository.getInvoices()
	data, _ := json.Marshal(invoices)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//NewInvoice POST
func (c *Controller) AddInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice Invoice
	//Read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Setting up Invoice: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Creating Invoice")
	}

	if err := json.Unmarshal(body, &invoice); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		success := c.repository.addInvoice(invoice)

		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	}
}

// UpdateInvoice PUT
func (c *Controller) ModifyInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice Invoice
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Invoice", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Invoice", err)
	}

	if err := json.Unmarshal(body, &invoice); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.modifyInvoice(invoice) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//FetchInvoice GET
func (c *Controller) FetchInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	invoice := c.repository.fetchInvoice(id)
	data, _ := json.Marshal(invoice)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//Delete an Invoice
func (c *Controller) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.deleteInvoice(id); err != "" {
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return

}

//Customers /customers GET
func (c *Controller) ListCustomer(w http.ResponseWriter, r *http.Request) {
	customers := c.repository.getCustomers()
	data, _ := json.Marshal(customers)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//NewCustomer POST
func (c *Controller) AddCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	//read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Setting Up Customer: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Creating Customer", err)
	}

	if err := json.Unmarshal(body, &customer); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		success := c.repository.addCustomer(customer)

		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	}
}

//FetchCustomer GET
func (c *Controller) FetchCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	customer := c.repository.fetchCustomer(id)
	data, _ := json.Marshal(customer)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// ModifyCustomer PUT
func (c *Controller) ModifyCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Customer", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Customer", err)
	}

	if err := json.Unmarshal(body, &customer); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.modifyCustomer(customer) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//Delete a Customer
func (c *Controller) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.deleteCustomer(id); err != "" {
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return

}

//Estimates /estimates GET
func (c *Controller) ListEstimate(w http.ResponseWriter, r *http.Request) {
	estimates := c.repository.getEstimates()
	data, _ := json.Marshal(estimates)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//NewEstimate POST
func (c *Controller) AddEstimate(w http.ResponseWriter, r *http.Request) {
	var estimate Estimate
	//read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Creating Estimate: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Creating Estimate", err)
	}
	if err := json.Unmarshal(body, &estimate); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		success := c.repository.AddEstimate(estimate)

		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	}
}

//FetchEstimate GET
func (c *Controller) FetchEstimate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	estimate := c.repository.fetchEstimate(id)
	data, _ := json.Marshal(estimate)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// UpdateEstimate PUT
func (c *Controller) ModifyEstimate(w http.ResponseWriter, r *http.Request) {
	var estimate Estimate
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Estimate", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Estimate", err)
	}

	if err := json.Unmarshal(body, &estimate); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.modifyEstimate(estimate) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//Delete an Estimate
func (c *Controller) DeleteEstimate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.deleteEstimate(id); err != "" {
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return

}

//Expenses /expenses GET
func (c *Controller) ListExpense(w http.ResponseWriter, r *http.Request) {
	expenses := c.repository.getExpenses()
	data, _ := json.Marshal(expenses)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//NewExpense POST
func (c *Controller) AddExpense(w http.ResponseWriter, r *http.Request) {
	var expense Expense
	//read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Creating Expense: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Creating Expense", err)
	}
	if err := json.Unmarshal(body, &expense); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		success := c.repository.AddExpense(expense)

		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	}
}

//FetchExpense GET
func (c *Controller) FetchExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	expense := c.repository.fetchExpense(id)
	data, _ := json.Marshal(expense)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// UpdateExpense PUT
func (c *Controller) ModifyExpense(w http.ResponseWriter, r *http.Request) {
	var expense Expense
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Expense", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Expense", err)
	}

	if err := json.Unmarshal(body, &expense); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.modifyExpense(expense) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//Delete an Expense
func (c *Controller) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.deleteExpense(id); err != "" {
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return

}

//Items /items GET
func (c *Controller) ListItem(w http.ResponseWriter, r *http.Request) {
	items := c.repository.getItems()
	data, _ := json.Marshal(items)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//NewItem POST
func (c *Controller) AddItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	//read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Creating Item: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Creating Item", err)
	}
	if err := json.Unmarshal(body, &item); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		success := c.repository.AddItem(item)

		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	}
}

// UpdateItem PUT
func (c *Controller) ModifyItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Item", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Item", err)
	}

	if err := json.Unmarshal(body, &item); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.modifyItem(item) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//FetchItem GET
func (c *Controller) FetchItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	item := c.repository.fetchItem(id)
	data, _ := json.Marshal(item)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//Delete an Item
func (c *Controller) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.deleteItem(id); err != "" {
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return

}

//Payments /payments GET
func (c *Controller) ListPayment(w http.ResponseWriter, r *http.Request) {
	invoices := c.repository.getInvoices()
	data, _ := json.Marshal(invoices)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//NewPayment POST
func (c *Controller) AcceptPayment(w http.ResponseWriter, r *http.Request) {
	var payment Payment
	//read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Creating Payment: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Creating Payment", err)
	}
	if err := json.Unmarshal(body, &payment); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		success := c.repository.AddPayment(payment)

		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	}
}

//FetchPayment GET
func (c *Controller) FetchPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	payment := c.repository.fetchPayment(id)
	data, _ := json.Marshal(payment)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// UpdatePayment PUT
func (c *Controller) ModifyPayment(w http.ResponseWriter, r *http.Request) {
	var payment Payment
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Payment", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Payment", err)
	}

	if err := json.Unmarshal(body, &payment); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.modifyPayment(payment) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//Delete a Payment
func (c *Controller) DeletePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.deletePayment(id); err != "" {
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return

}

//Items /taxes GET
func (c *Controller) ListTax(w http.ResponseWriter, r *http.Request) {
	taxes := c.repository.getTaxes()
	data, _ := json.Marshal(taxes)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//NewTax POST
func (c *Controller) AddTax(w http.ResponseWriter, r *http.Request) {
	var tax Tax
	//read the body of the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Creating Item: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Creating Item", err)
	}
	if err := json.Unmarshal(body, &tax); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		success := c.repository.AddTax(tax)

		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		return
	}
}

//FetchTax GET
func (c *Controller) FetchTax(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	tax := c.repository.fetchTax(id)
	data, _ := json.Marshal(tax)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// UpdateTax PUT
func (c *Controller) ModifyTax(w http.ResponseWriter, r *http.Request) {
	var tax Tax
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error Updating Tax", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Updating Tax", err)
	}

	if err := json.Unmarshal(body, &tax); err != nil {
		// unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.repository.modifyTax(tax) // updates the product in the DB

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

//Delete a Tax
func (c *Controller) DeleteTax(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.deleteTax(id); err != "" {
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return

}
func (c *Controller) ConfirmSubscription(w http.ResponseWriter, r *http.Request) {

}
func (c *Controller) SendEstimate(w http.ResponseWriter, r *http.Request) {

}
func (c *Controller) SendInvoice(w http.ResponseWriter, r *http.Request) {

}
