package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Product struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price float64
	Stock int
}

type Order struct {
	ID       uint `gorm:"primaryKey"`
	Products []ProductOrder
}

type ProductOrder struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint
	Product   Product
	OrderID   uint
	Order     Order
	Amount    int
}

type User struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Profile  StudentProfile
}

type StudentProfile struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	CompanyName string
	JobTitle    string
	Level       string
}

type Course struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
}

type Class struct {
	ID        uint `gorm:"primaryKey"`
	CourseID  uint
	Course    Course
	TrainerID uint
	Trainer   User
	Start     time.Time
	End       time.Time
	Seats     int
	Students  []ClassStudent
}

type ClassStudent struct {
	ID        uint `gorm:"primaryKey"`
	ClassID   uint
	StudentID uint
	Student   User
}

func main() {
	url := "host=localhost user=peagolang password=supersecret dbname=peagolang port=54329 sslmode=disable"
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(
		&Product{},
		&Order{},
		&User{},
		&StudentProfile{},
		&ProductOrder{},
		&Course{},
		&Class{},
		&ClassStudent{},
	)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().AutoMigrate(
		&Product{},
		&Order{},
		&User{},
		&StudentProfile{},
		&ProductOrder{},
		&Course{},
		&Class{},
		&ClassStudent{},
	)
	if err != nil {
		log.Fatal(err)
	}

	tdd := Course{
		Name:        "TDD",
		Description: "TDD is fun!",
	}
	db.Create(&tdd)

	pong := User{Username: "pong"}
	gap := User{Username: "gap"}
	kane := User{Username: "kane"}
	jua := User{Username: "jua"} // Trainer

	db.Create(&pong)
	db.Create(&gap)
	db.Create(&kane)
	db.Create(&jua)

	class := Class{
		CourseID:  tdd.ID,
		TrainerID: jua.ID,
		Start:     time.Date(2023, 5, 10, 9, 0, 0, 0, time.Local),
		End:       time.Date(2023, 5, 12, 17, 0, 0, 0, time.Local),
		Seats:     10,
		Students: []ClassStudent{
			{StudentID: pong.ID},
			{StudentID: gap.ID},
		},
	}

	db.Save(&class)

	var foundClass Class
	db.Preload("Course").Preload("Trainer").Preload("Students.Student").First(&foundClass, class.ID)

	// select * from class cls
	// join course cou on cou.id = cls.class_id
	// join user trainer on trainer.id = cls.user_id
	// join class_students clst on clst.id = cls.class_student_id
	// join user student on student.id = clst.student_id
	// where id = 1;

	fmt.Println("#ID: ", foundClass.ID)
	fmt.Println("Name: ", foundClass.Course.Name)
	fmt.Println("Description: ", foundClass.Course.Description)
	fmt.Println("\tBy: ", foundClass.Trainer.Username)
	fmt.Println("\tDate: ", foundClass.Start, foundClass.End)
	fmt.Println("Students: ")
	for _, student := range foundClass.Students {
		fmt.Println("\tName: ", student.Student.Username)
	}

	// shirt := Product{
	// 	Name:  "T-Shirt",
	// 	Price: 350,
	// 	Stock: 200,
	// }

	// short := Product{
	// 	Name:  "Short v1",
	// 	Price: 600,
	// 	Stock: 150,
	// }

	// toy := Product{
	// 	Name:  "Car Toy",
	// 	Price: 99,
	// 	Stock: 700,
	// }

	// db.Create(&shirt)
	// db.Create(&short)
	// db.Create(&toy)

	// order1 := Order{
	// 	Products: []ProductOrder{
	// 		{ProductID: shirt.ID, Amount: 1},
	// 		{ProductID: short.ID, Amount: 1},
	// 	},
	// }
	// db.Create(&order1)

	// order2 := Order{
	// 	Products: []ProductOrder{
	// 		{ProductID: shirt.ID, Amount: 1},
	// 		{ProductID: toy.ID, Amount: 1},
	// 	},
	// }
	// db.Create(&order2)

	// var foundOrder Order
	// db.Preload("Products.Product").First(&foundOrder, order1.ID)
	// fmt.Printf("\n\n%+v\n\n", foundOrder)
	// PrintOrder(foundOrder)
}

// func PrintOrder(order Order) {
// 	fmt.Println()
// 	fmt.Printf("Order ID: %v\n", order.ID)
// 	fmt.Println("Products:")
// 	for _, p := range order.Products {
// 		fmt.Printf("\t%v\t\t%v\t%v\n", p.Product.Name, p.Product.Price, p.Amount)
// 	}
// }

// var found Product
// db.Preload("Orders").First(&found, 1)

// fmt.Printf("\n\n %+v \n\n", found)

// var found2 Order
// db.Preload("Product").First(&found2, 1)

// fmt.Printf("\n\n %+v \n\n", found2)

// user := User{
// 	Username: "pong",
// 	Profile: StudentProfile{
// 		CompanyName: "ODDS",
// 		JobTitle:    "Golang Developer",
// 		Level:       "Poring",
// 	},
// }

// db.Save(&user)

// var foundUser User
// db.Preload("Profile").First(&foundUser, user.ID)

// fmt.Println(foundUser)
