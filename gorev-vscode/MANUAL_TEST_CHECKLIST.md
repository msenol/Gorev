# Subtask UI Manual Test Checklist

## 🚀 Başlangıç
- [ ] VS Code'u F5 ile debug modda başlat
- [ ] MCP Server'ı `./gorev serve --debug` ile başlat
- [ ] Gorev extension'ına bağlan (Connect butonu)
- [ ] Test için bir proje oluştur

## 📝 Alt Görev Oluşturma

### Sağ Tık Menüsü ile
- [ ] Bir görev oluştur
- [ ] Görev üzerine sağ tıkla
- [ ] "Create Subtask" seçeneğini gör
- [ ] Alt görev bilgilerini gir (başlık, açıklama, öncelik)
- [ ] Alt görevin oluşturulduğunu doğrula
- [ ] Parent görevin yanında genişletme oku olduğunu doğrula
- [ ] Parent görevi genişlet ve alt görevi gör

### Drag & Drop ile
- [ ] İki bağımsız görev oluştur
- [ ] Bir görevi diğerinin üzerine sürükle
- [ ] "Alt Görev Yap" ve "Bağımlılık Oluştur" seçeneklerini gör
- [ ] "Alt Görev Yap" seçeneğini seç
- [ ] Alt görevin oluşturulduğunu doğrula

## 🔄 Parent Değiştirme

### Sağ Tık ile
- [ ] Alt görev üzerine sağ tıkla
- [ ] "Change Parent Task" seçeneğini gör
- [ ] Görev listesinden yeni parent seç
- [ ] Parent'ın değiştiğini doğrula

### Drag & Drop ile
- [ ] Alt görevi başka bir görevin üzerine sürükle
- [ ] "Alt Görev Yap" seçeneğini seç
- [ ] Parent'ın değiştiğini doğrula

## 🚫 Parent Kaldırma

### Sağ Tık ile
- [ ] Alt görev üzerine sağ tıkla
- [ ] "Remove Parent (Make Root Task)" seçeneğini gör
- [ ] Seçeneği tıkla
- [ ] Görevin artık root level'da göründüğünü doğrula

### Drag & Drop ile
- [ ] Alt görevi boş alana sürükle
- [ ] Görevin root level'a taşındığını doğrula

## 🎯 Hiyerarşik Görüntüleme

### TreeView
- [ ] Parent görevlerin yanında genişletme oku var
- [ ] Alt görev sayısı gösteriliyor (📁 2/5 formatında)
- [ ] Tamamlanan alt görev sayısı doğru
- [ ] Alt görevler indent edilmiş şekilde gösteriliyor
- [ ] Çoklu seviye hiyerarşi düzgün gösteriliyor

### Task Detail Panel
- [ ] Parent göreve tıkla
- [ ] Hiyerarşi bölümü görünüyor
- [ ] Toplam alt görev sayısı doğru
- [ ] İlerleme yüzdesi doğru hesaplanmış
- [ ] İlerleme çubuğu doğru oranda dolu
- [ ] "Alt Görev Oluştur" butonu çalışıyor

## ⚠️ Hata Senaryoları

### Dairesel Bağımlılık
- [ ] A görevini B'nin altına taşı
- [ ] B görevini A'nın altına taşımayı dene
- [ ] "Dairesel bağımlılık" hatası gösteriliyor

### Farklı Proje Kısıtlaması
- [ ] İki farklı proje oluştur
- [ ] Proje 1'de bir görev oluştur
- [ ] Proje 2'de bir görev oluştur
- [ ] Bir görevi diğer projedeki görevin altına taşımayı dene
- [ ] "Aynı projede olmalı" hatası gösteriliyor

## 🎨 UI/UX Kontrolleri

### Context Values
- [ ] Root görevlerde context menü öğeleri doğru
- [ ] Parent görevlerde "task:parent" context value
- [ ] Child görevlerde "task:child" context value
- [ ] Child görevlerde "Remove Parent" seçeneği var
- [ ] Tüm görevlerde "Create Subtask" seçeneği var

### Görsel İndikatörler
- [ ] Parent görevler farklı ikon gösteriyor
- [ ] Alt görev sayısı badge'i görünüyor
- [ ] Genişletme/daraltma animasyonu çalışıyor
- [ ] Drag & drop sırasında görsel feedback var

## 🔧 Konfigürasyon

### Ayarları Test Et
- [ ] Settings > Gorev > Drag Drop > Allow Parent Change ayarını kapat
- [ ] Drag & drop ile parent değiştirmenin devre dışı olduğunu doğrula
- [ ] Ayarı tekrar aç ve çalıştığını doğrula

## 📊 Performans

### Büyük Hiyerarşiler
- [ ] 10+ alt görevi olan bir parent oluştur
- [ ] 3+ seviye derinliğinde hiyerarşi oluştur
- [ ] TreeView'ın hızlı yüklendiğini doğrula
- [ ] Genişletme/daraltmanın hızlı olduğunu doğrula

## 🐛 Bilinen Sorunlar
- [ ] Çok hızlı drag & drop işlemlerinde UI güncellemesi gecikebilir
- [ ] 100+ alt görevde performans düşebilir

## ✅ Test Tamamlama
- [ ] Tüm temel fonksiyonlar çalışıyor
- [ ] Hata senaryoları düzgün ele alınıyor
- [ ] UI güncellemeleri doğru yapılıyor
- [ ] Performans kabul edilebilir seviyede

---

Test Tarihi: _______________
Test Eden: _______________
Versiyon: 0.8.0