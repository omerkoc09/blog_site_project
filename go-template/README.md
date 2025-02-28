## Yazılım Projesi Prensipleri
- Bir projenin en önemli kısmı veridir. (Misal: bir e-ticaret sitesi için ürünler(ismi, boyu, kilosu, rengi vb.), kategoriler(elektronik, beyazeşya vb.), kullanıcılar (isim, soyisim, adres vb.), vb.)
- Projeye veriyi modelleyerek başlanır.
- Veri modellemesi yapılırken verinin nasıl görüneceği değil, verinin nasıl tutulacağı düşünülür.
- Projedeki 2. önemli kısım Business Logic'dir. Veriyi nasıl işleyeceğimiz. Hangi kurallar dahilinde veri tutulacak veya gösterileceğinin belirlenmesidir.
- 

## Yazılım Projesi Veri Modellemesi
- Veri veritabanında tutulur. (Misal: biz postgresqlde tutuyoruz)
- Veritabanı bizim projemiz haricinde ayrı bir yazılımdır. postgres ayrı olarak çalışır ve bir port ile ona bağlanıp sql sorguları çalıştırırız.
- Veriyi go tarafında modeller ile temsil ederiz. (Misal: bir ürün modeli, bir kategori modeli, bir kullanıcı modeli vb.)
- Model go'da structa karşılık gelir. Veritabanındaki tablo = model, kolon = struct fieldları.
- Modellerin veritabanındaki karşılığını olşuturma işlemine migration denir.
- Go'da istisnalar hariç sql sorgusu yazmıyoruz, ORM ile sorguları oluşturuyoruz.

## Go Proje Yapımız
- Backend uygulamalarımız backend/ dizini altında tutulur. go.mod dosyası da bu dizinde bulunur.
- pkg/ de projenin genelinde kullanılacak paketler tutulur. (Misal: logger, error, db bağlantısı, base handler/service, utils vb.)
  - pkg/app de Haytek web frameworkü bulunuyor.
  - pkg/cmd de proje için veritabanı ortak işlemleri bulunur. (Misal: migrate, rollback, seed vb.)
  - pkg/config de projeye göre özelleştirilebilecek ayarlar bulunur.
  - pkg/database de veritabanı bağlantısı ve işlemleri bulunur.
  - pkg/errorsx de Haytek frameworküne özgü hata yönetimi bulunur.
  - pkg/

## Go Projesi Örnek Api implementasyonu
- backend/common/model içinde veritabanı modeli oluşturulur. IDBModel interface'ini implement eder.
- backend/common/migrations/20220902105002_bismillah/ veri tabanı modeli kopyalanır. Başlangıçta model/ klasörü ile bismillah klasörü benzer içeriği taşır. Ancak proje ilerledikçe, canlı ortamda değişikliker yapıldıkça migrations altına yeni dosyalar eklenir.
- model/ klasörü içinde veritabanının son hali bulunur. migrations/ altındaki dosyalar ise veritabanının geçmiş hallerini temsil eder.
- `make migrate` komutu ile migrations/ altındaki dosyaları veritabanına uygularız.
- `make seed` komutu ile common/fixture/fixture.yml dosyasındaki verileri veritabanına yükleriz.
- **DEVELOPMENT** esnasında **YANİ CANLI DA DEĞİLKEN**, **GERÇEK VERİ YOKKEN** modelde bir değişlik yaptık. veritabanını silip oluşturabiliriz. bunun için  `make resetdb` komutunu kullanırız. 
  - Bu komut veritabanını siler ve tekrar oluşturur. Bu işlemi yaparken veritabanındaki veriler silinir. Bu yüzden dikkatli kullanılmalıdır.
- Canlıya çıkasıya kadar migartions/bismillahın altına modeller eklenir ve resetdb yapılır.
- Diyelim ki canlıya çıktık. Gerçek veri var. Artık resetdb yapamayız. Bu durumda migrations/ altına yeni bir dosya ekleyip migration yapılır.
  - birden fazla model oluşturulacaksa `make migrationdir name=migration_klasor_ismi` komutu ile migrations/ altına yeni bir klasör oluşturulur. ve modeller eklendikten sonra `make migrate` ile veritabanına uygulanır.
  - tek model oluşturulacaksa `make migrationfile name=migration_dosya_ismi` komutu ile migrations/ altına yeni bir dosya oluşturulur. ve modeller eklendikten sonra `make migrate` ile veritabanına uygulanır.
- service klasörünün altına ilgili model için service oluşturulur. (Misal: product service, category service, user service vb.)
- backend/idare/viewmodel ile ilgili modelin viewmodeli oluşturulur. (Misal: product viewmodel, category viewmodel, user viewmodel vb.)
- backend/idare/handler ile ilgili modelin handlerı oluşturulur. (Misal: product handler, category handler, user handler vb.)
- en sonda backend/idare/router.go içinde ilgili handlerlar route edilir.


## Best Practices
- Veritabanı işlemleri servis içinde yapılır. Handlerda veritabanı işlemi yapılmaz.
- Businnes Logicler servislerde tutulur. Handlerda business logic yazılmaz!


## Gotemplatede migration
- local postgres sunucusu için `make devdb` // localde kurulu postgres ya onu durdur veya portları değiştirebilirsin. bu 5432de çalışıyor
- veritabanı init: `make migrate` common/migrations/20220902105002_bismillah.go da bulunan modellerin tablolarını oluşturur ve 
- veritabanı ilk veriler: `make seed` common/fixture/fixture.yml içinde bulunan verileri dbye yükler. bu komu her çalıştığında ilgili tablolardaki verililer komple silinir yeiden olşuturulur
- 

