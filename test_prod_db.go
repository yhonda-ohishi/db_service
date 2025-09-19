package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/repository"
)

func main() {
	// 環境変数ロード
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file could not be loaded: %v", err)
	}

	// 本番データベース接続テスト
	fmt.Println("本番データベース接続テスト開始...")

	prodDB, err := config.NewProdDatabase()
	if err != nil {
		log.Fatalf("本番DB接続エラー: %v", err)
	}
	defer prodDB.Close()

	fmt.Println("本番DB接続成功!")

	// DTakoCarsテーブルテスト
	fmt.Println("\n=== DTakoCars テスト ===")
	carsRepo := repository.NewDTakoCarsRepository(prodDB)

	cars, totalCount, err := carsRepo.GetAll(5, 0)
	if err != nil {
		log.Printf("DTakoCars取得エラー: %v", err)
	} else {
		fmt.Printf("車輌データ総数: %d件\n", totalCount)
		for i, car := range cars {
			fmt.Printf("%d. ID:%d, 車輌CD:%s, 車輌名:%s\n", i+1, car.ID, car.CarCode, car.CarName)
		}
	}

	// DTakoEventsテーブルテスト
	fmt.Println("\n=== DTakoEvents テスト ===")
	eventsRepo := repository.NewDTakoEventsRepository(prodDB)

	events, totalCount, err := eventsRepo.GetAll(3, 0)
	if err != nil {
		log.Printf("DTakoEvents取得エラー: %v", err)
	} else {
		fmt.Printf("イベントデータ総数: %d件\n", totalCount)
		for i, event := range events {
			fmt.Printf("%d. ID:%d, 運行NO:%s, イベント名:%s\n", i+1, event.ID, event.OperationNo, event.EventName)
		}
	}

	// DTakoRowsテーブルテスト
	fmt.Println("\n=== DTakoRows テスト ===")
	rowsRepo := repository.NewDTakoRowsRepository(prodDB)

	rows, totalCount, err := rowsRepo.GetAll(3, 0)
	if err != nil {
		log.Printf("DTakoRows取得エラー: %v", err)
	} else {
		fmt.Printf("運行データ総数: %d件\n", totalCount)
		for i, row := range rows {
			fmt.Printf("%d. ID:%s, 運行NO:%s, 車輌CD:%d\n", i+1, row.ID, row.OperationNo, row.CarCode)
		}
	}

	// ETCNumテーブルテスト
	fmt.Println("\n=== ETCNum テスト ===")
	etcRepo := repository.NewETCNumRepository(prodDB)

	etcNums, totalCount, err := etcRepo.GetAll(5, 0)
	if err != nil {
		log.Printf("ETCNum取得エラー: %v", err)
	} else {
		fmt.Printf("ETCカードデータ総数: %d件\n", totalCount)
		for i, etc := range etcNums {
			fmt.Printf("%d. ETCカード番号:%s, 車輌ID:%s\n", i+1, etc.ETCCardNum, etc.CarID)
		}
	}

	fmt.Println("\n本番DB接続テスト完了!")
}
