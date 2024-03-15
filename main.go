package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

type BookLibrary struct {
	KodeBuku      string
	JudulBuku     string
	Pengarang     string
	Penerbit      string
	JumlahHalaman int
	TahunTerbit   int
	Tanggal       time.Time
}

var ListBook []BookLibrary

func tambahBuku() {
	fmt.Println("\n")
	bookCode := ""
	bookTitle := ""
	bookAuthor := ""
	bookPublisher := ""
	var pageTotal int
	var publishedYear int

	fmt.Println("===================")
	fmt.Println("Tambah Buku")
	fmt.Println("===================")
	fmt.Println(ListBook)

	draftBuku := []BookLibrary{}

	for {
		fmt.Print("Masukkan Code Buku : ")
		_, err := fmt.Scanln(&bookCode)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
	
		for _, book := range draftBuku {
			if book.KodeBuku == "book-"+bookCode {
				fmt.Println("Kode Buku Sudah Ada. Masukkan Kode Buku Lain.")
				return
			}
		}
	
		fmt.Print("Masukkan Judul Buku : ")
		_, err = fmt.Scanln(&bookTitle)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
	
		fmt.Print("Masukkan Pengarang Buku : ")
		_, err = fmt.Scanln(&bookAuthor)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
	
		fmt.Print("Masukkan Penerbit Buku : ")
		_, err = fmt.Scanln(&bookPublisher)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
	
		fmt.Print("Masukkan Total Halaman : ")
		_, err = fmt.Scanln(&pageTotal)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
	
		fmt.Print("Masukkan Tahun Terbit : ")
		_, err = fmt.Scanln(&publishedYear)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}

		draftBuku = append(draftBuku, BookLibrary{
			KodeBuku: fmt.Sprintf("book-%s", bookCode),
			JudulBuku: bookTitle,
			Pengarang: bookAuthor,
			Penerbit: bookPublisher,
			JumlahHalaman: pageTotal,
			TahunTerbit: publishedYear,
			Tanggal: time.Now(),
		})

		pilihanTambahBuku := 0
		fmt.Println("Ketik 1 untuk tambah buku, ketik 0 untuk keluar")
		_, err = fmt.Scanln(&pilihanTambahBuku)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		if pilihanTambahBuku == 0 {
			break
		}
	}

	fmt.Println("Menambah Buku...")

	_ = os.Mkdir("buku", 0777)

	ch := make(chan BookLibrary)

	wg := sync.WaitGroup{}

	jumlahPustakawan := 5

	// Mendaftarkan receiver/pemroses data
	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go simpanBuku(ch, &wg, i)
	}

	for _, bukuTersimpan := range draftBuku {
		ch <- bukuTersimpan
	}

	close(ch)

	wg.Wait()

	fmt.Println("Berhasil Tambah Buku")
}

func simpanBuku(ch <-chan BookLibrary, wg *sync.WaitGroup, noPustakawan int)  {

	for bukuTersimpan := range ch {
		dataJson, err := json.Marshal(bukuTersimpan)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		err = os.WriteFile(fmt.Sprintf("buku/%s.json", bukuTersimpan.KodeBuku), dataJson, 0644)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		fmt.Printf("Pustakawan No %d Memproses Kode Buku : %s!\n", noPustakawan, bukuTersimpan.KodeBuku)
	}
	wg.Done()
}

func lihatListBuku(ch <-chan string, chBuku chan BookLibrary, wg *sync.WaitGroup)  {
	var bookLibrary BookLibrary
	for kodeBuku := range ch {
		dataJson, err := os.ReadFile(fmt.Sprintf("buku/%s", kodeBuku))
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		err = json.Unmarshal(dataJson, &bookLibrary)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		chBuku <- bookLibrary
	}
	wg.Done()
}

func listBuku() {
	fmt.Println("\n")
	fmt.Println("===================")
	fmt.Println("List Buku")
	fmt.Println("===================")
	fmt.Println("Memuat Data...")
	ListBook = []BookLibrary{}

	listJsonBuku,err :=  os.ReadDir("buku")
	if err != nil {
		fmt.Println("Terjadi error: ", err)
	}

	wg := sync.WaitGroup{}

	ch := make(chan string)
	chBuku := make(chan BookLibrary, len(listJsonBuku))

	jumlahPustakawan := 5

	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go lihatListBuku(ch, chBuku, &wg)
	}

	for _, fileBuku := range listJsonBuku {
		ch <- fileBuku.Name()
	}

	close(ch)

	wg.Wait()

	close(chBuku)

	for dataBookLibrary := range chBuku {
		ListBook = append(ListBook, dataBookLibrary)
	}

	sort.Slice(ListBook, func(i, j int) bool {
		return ListBook[i].Tanggal.Before(ListBook[j].Tanggal)
	})

	if len(ListBook) < 1 {
		fmt.Println("---Tidak Ada Buku---")
	}

	for i, v := range ListBook {
		i++
		fmt.Printf("%d. Kode Buku : %s, Judul Buku : %s, Pengarang : %s, Penerbit : %s, Jumlah Halaman : %d, Tahun Terbit : %d\n", i, v.KodeBuku, v.JudulBuku, v.Pengarang, v.Penerbit, v.JumlahHalaman, v.TahunTerbit)
	}
}

func detailBuku(kode string) {
	fmt.Println("\n")
	fmt.Println("===================")
	fmt.Println("Detail Buku")
	fmt.Println("===================")

	var isBook bool

	for _, book := range ListBook {
		if book.KodeBuku == kode {
			isBook = true
			fmt.Printf("Kode Buku : %s\n", book.KodeBuku)
			fmt.Printf("Judul Buku : %s\n", book.JudulBuku)
			fmt.Printf("Pengarang Buku : %s\n", book.Pengarang)
			fmt.Printf("Penerbit Buku : %s\n", book.Penerbit)
			fmt.Printf("Jumlah Halaman : %d\n", book.JumlahHalaman)
			fmt.Printf("Tahun Terbit : %d\n", book.TahunTerbit)
			break
		}
	}

	if !isBook {
		fmt.Println("Kode Buku Salah Atau Tidak Ada")
	}
}

func updateBuku(kode string) {
	fmt.Println("\n")
	detailBuku(kode)

	fmt.Println("===================")
	fmt.Println("Edit Buku")
	fmt.Println("===================")

	var book BookLibrary

	fmt.Print("Masukkan Code Buku : ")
	_, err := fmt.Scanln(&book.KodeBuku)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Judul Buku : ")
	_, err = fmt.Scanln(&book.JudulBuku)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Pengarang Buku : ")
	_, err = fmt.Scanln(&book.Pengarang)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Penerbit Buku : ")
	_, err = fmt.Scanln(&book.Penerbit)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Total Halaman : ")
	_, err = fmt.Scanln(&book.JumlahHalaman)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Tahun Terbit : ")
	_, err = fmt.Scanln(&book.TahunTerbit)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Println(book)

	for i, b := range ListBook {
		if b.KodeBuku == kode {
			ListBook[i] = book
			break
		}
	}
}

func hapusBuku(kode string) {
	fmt.Println("\n")
	var isBook bool
	for i, book := range ListBook {
		if book.KodeBuku == kode {
			isBook = true
			err := os.Remove(fmt.Sprintf("buku/%s.json", ListBook[i].KodeBuku))
			if err != nil {
				fmt.Println("Terjadi error:", err)
			}
			fmt.Println("Buku Berhasil Dihapus")
			break
		}
	}


	if !isBook {
		fmt.Println("Kode Buku Salah Atau Tidak Ada")
	}
}

func main() {
	fmt.Println("\n")

	var opsi int

	fmt.Println("===================")
	fmt.Println("Manajemen Buku Perpustakaan")
	fmt.Println("===================")

	fmt.Println("Pilih Opsi")
	fmt.Println("1. Tambah Buku")
	fmt.Println("2. Lihat Semua Buku")
	fmt.Println("3. Lihat Detail Buku")
	fmt.Println("4. Edit Buku")
	fmt.Println("5. Hapus Buku")
	fmt.Println("6. Keluar")

	fmt.Print("Masukkan Opsi : ")
	_, err := fmt.Scanln(&opsi)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	switch opsi {
	case 1:
		tambahBuku()
	case 2:
		listBuku()
	case 3:
		var pilihDetail string
		listBuku()
		fmt.Print("Masukkan Kode Buku : ")
		_, err := fmt.Scanln(&pilihDetail)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
		detailBuku(pilihDetail)
	case 4:
		var pilihanUpdate string
		listBuku()
		fmt.Print("Masukkan Kode Buku Yang Akan DiEdit : ")
		_, err := fmt.Scanln(&pilihanUpdate)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
		updateBuku(pilihanUpdate)
	case 5:
		var pilihanHapus string
		listBuku()
		fmt.Print("Masukkan Kode Buku Yang Akan Dihapus : ")
		_, err := fmt.Scanln(&pilihanHapus)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
		hapusBuku(pilihanHapus)
	case 6:
		os.Exit(0)
	default:
		fmt.Println("Tidak Ada Opsi")
	}

	main()
}