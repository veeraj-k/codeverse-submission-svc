package db

import (
	"fmt"
	"os"
	"submission-service/internal/submission"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return DB
}

func MigrateDB() {
	DB.Migrator().DropTable(&submission.Submission{})
	DB.AutoMigrate(&submission.Submission{})
}

func SeedDB() {
	submissions := []submission.Submission{
		{
			UserId:          1,
			ProblemId:       101,
			Language:        "Go",
			Code:            "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			Status:          "Accepted",
			TestCasesPassed: 10,
			TotalTestCases:  10,
		},
		{
			UserId:          2,
			ProblemId:       102,
			Language:        "Python",
			Code:            "print('Hello, World!')",
			Status:          "Accepted",
			TestCasesPassed: 8,
			TotalTestCases:  10,
		},
		{
			UserId:          3,
			ProblemId:       103,
			Language:        "Java",
			Code:            "public class Main {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello, World!\");\n\t}\n}",
			Status:          "Wrong Answer",
			TestCasesPassed: 5,
			TotalTestCases:  10,
		},
		{
			UserId:          4,
			ProblemId:       104,
			Language:        "JavaScript",
			Code:            "console.log('Hello, World!');",
			Status:          "Accepted",
			TestCasesPassed: 10,
			TotalTestCases:  10,
		},
		{
			UserId:          5,
			ProblemId:       105,
			Language:        "C++",
			Code:            "#include <iostream>\n\nint main() {\n\tstd::cout << \"Hello, World!\" << std::endl;\n\treturn 0;\n}",
			Status:          "Compilation Error",
			TestCasesPassed: 0,
			TotalTestCases:  10,
		},
	}

	for _, s := range submissions {
		DB.Create(&s)
	}
}
