package sqlitecreator

import (
	_ "embed"
	"database/sql"
	"fmt"
	"log"
	"time"
	"bufio"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed money.txt
var moneyTxt string //-----------------moneyTxt変数にtxtを埋め込む


func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// コマンドメソッド
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		amount INTEGER,
		expenseCategory TEXT
	)`)
	if err != nil {
		return err
	}
	return nil
}


// コマンドメソッド
func ProcessEmbeddedFile(db *sql.DB, expenseCategoryIndex int) error {
    // 文字列として埋め込まれたファイルからScannerを作成
    scanner := bufio.NewScanner(strings.NewReader(moneyTxt))

    for scanner.Scan() {
        data := strings.Split(scanner.Text(), ",")

		if len(data) < 5 {
			log.Println("Invalid data format")
			continue
		}

		// 日付の形式をチェック
		if !isDateFormatValid(data[1]) {
			log.Println("無効な日付形式です:", data[1])
			continue
		}

		// nullチェック ExpenseCategory、デフォルト値として"dontworry"を設定
		expenseCategory := data[expenseCategoryIndex]
		if expenseCategory == "" {
			expenseCategory = "dontworry"
		}

		// データベースに挿入
		_, err := db.Exec("INSERT INTO entries (date, amount, expenseCategory) VALUES (?, ?, ?)",
			data[1], data[2], expenseCategory)
		if err != nil {
			log.Println("Error inserting into database:", err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// 日付の形式が有効かどうかを検証する関数
func isDateFormatValid(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}


//---処理型(DBカラム用)
type Entry struct {
	ID              int //割り振り直す
	Date            time.Time
	Amount          int
	ExpenseCategory string //クエリで使うカテゴリー
}

// データをクエリして構造体に変換し、標準出力するテスト(頭の5つのみ)
func OutputEntries(db *sql.DB) error {
	rows, err := db.Query("SELECT id, date, amount, expenseCategory FROM entries LIMIT 5")
	if err != nil {
		return err
	}
	defer rows.Close()

	var entry Entry
	for rows.Next() {
		var dateStr string // 文字列としての日付
		err := rows.Scan(&entry.ID, &dateStr, &entry.Amount, &entry.ExpenseCategory)
		if err != nil {
			return err
		}

		// 文字列を time.Time に変換
		entry.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", entry)
	}
	fmt.Println("-------確認終了---------")

	return nil
}

// クエリメソッド すべてのAmountの合計をSQLで計算して標準出力
func CalculateTotalAmount(db *sql.DB) error {
	rows, err := db.Query("SELECT SUM(amount) AS TotalAmount FROM entries")
	if err != nil {
		return err
	}
	defer rows.Close()

	var totalAmount int
	for rows.Next() {
		err := rows.Scan(&totalAmount)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Total Amount: %d\n", totalAmount)
	fmt.Println("-------確認終了---------")

	return nil
}


//---処理型(合計値出力用)
type MonthTotal struct {
	Month         string
	TotalAmount   int
}

// クエリメソッド 月ごとのAmount合計を計算
func CalculateTotalAmountPerMonth(db *sql.DB) ([]MonthTotal, error) {
	rows, err := db.Query("SELECT strftime('%Y-%m', date) AS month, SUM(amount) AS TotalAmount FROM entries GROUP BY month")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthTotals []MonthTotal

	for rows.Next() {
		var monthTotal MonthTotal
		if err := rows.Scan(&monthTotal.Month, &monthTotal.TotalAmount); err != nil {
			return nil, err
		}
		monthTotals = append(monthTotals, monthTotal)
	}

	return monthTotals, nil
}

func CalculateAverageTotalAmountPerMonth(monthTotals []MonthTotal) float64 {
	totalMonths := len(monthTotals)
	if totalMonths == 0 {
		return 0
	}

	totalAmount := 0
	for _, mt := range monthTotals {
		totalAmount += mt.TotalAmount
	}

	average := float64(totalAmount) / float64(totalMonths)
	return average
}


// クエリメソッド 
func CalculateTotalAmountPerMonthExcludingYatin(db *sql.DB) ([]MonthTotal, error) {
	rows, err := db.Query("SELECT strftime('%Y-%m', date) AS month, SUM(amount) AS TotalAmount FROM entries WHERE expenseCategory != 'yatin' GROUP BY month")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totals []MonthTotal

	for rows.Next() {
		var total MonthTotal
		if err := rows.Scan(&total.Month, &total.TotalAmount); err != nil {
			return nil, err
		}
		totals = append(totals, total)
	}

	return totals, nil
}

func CalculateAverageTotalAmountPerMonthExcludingYatin(totals []MonthTotal) float64 {
	var sum, count int
	for _, t := range totals {
		sum += t.TotalAmount
		count++
	}

	if count == 0 {
		return 0
	}

	average := float64(sum) / float64(count)
	return average
}


func Sqlitecreator() {
	db, err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	// txtのIndex=4をexpenseCategoryとして使用する
    err = ProcessEmbeddedFile(db, 4)
    if err != nil {
        log.Fatal(err)
    }

	err = OutputEntries(db)
	if err != nil {
		log.Fatal(err)
	}

	err = CalculateTotalAmount(db)
	if err != nil {
		log.Fatal(err)
	}

	// 月ごとのAmount合計を計算
	monthTotals, err := CalculateTotalAmountPerMonth(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, mt := range monthTotals {
		fmt.Printf("Total Amount for %s: %d\n", mt.Month, mt.TotalAmount)
	}
	average := CalculateAverageTotalAmountPerMonth(monthTotals)
	fmt.Printf("Average Total Amount per Month: %.2f\n", average)
	fmt.Println("-------確認終了---------")

	// 月ごとのAmount合計を計算（yatinのAmountを除外）
	totals, err := CalculateTotalAmountPerMonthExcludingYatin(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range totals {
		fmt.Printf("Total Amount for %s (Excluding yatin): %d\n", t.Month, t.TotalAmount)
	}
	average2 := CalculateAverageTotalAmountPerMonthExcludingYatin(totals)
	fmt.Printf("Average Total Amount per Month (Excluding yatin): %.2f\n", average2)
	fmt.Println("-------確認終了---------")
}