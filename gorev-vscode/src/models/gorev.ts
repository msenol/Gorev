import { GorevDurum, GorevOncelik, Timestamp } from './common';

export { GorevDurum, GorevOncelik } from './common';

export interface Gorev extends Timestamp {
  id: string;
  baslik: string;
  aciklama: string;
  durum: 'beklemede' | 'devam_ediyor' | 'tamamlandi';
  oncelik: 'dusuk' | 'orta' | 'yuksek';
  proje_id?: string;
  parent_id?: string;
  son_tarih?: string;
  etiketler?: Array<{ id: string; isim: string }>;
  proje_name?: string;
  bagimliliklar?: Bagimlilik[];
  alt_gorevler?: Gorev[];
  seviye?: number;
  // Dependency count fields for TreeView display
  bagimli_gorev_sayisi?: number; // Number of dependencies this task has
  tamamlanmamis_bagimlilik_sayisi?: number; // Number of incomplete dependencies
  bu_goreve_bagimli_sayisi?: number; // Number of tasks that depend on this task
}

export interface GorevDetay extends Gorev {
  proje_isim: string;
  bagimliliklar?: Bagimlilik[];
}

export interface GorevHiyerarsi {
  gorev: Gorev;
  ust_gorevler: Gorev[];
  toplam_alt_gorev: number;
  tamamlanan_alt: number;
  devam_eden_alt: number;
  beklemede_alt: number;
  ilerleme_yuzdesi: number;
}

export interface Bagimlilik {
  kaynak_id: string;
  hedef_id: string;
  baglanti_tip: string;
  hedef_baslik?: string;
  hedef_durum?: GorevDurum;
}

export interface BagimlilikOzet {
  toplam_bagimlilik: number;
  tamamlanan_bagimlilik: number;
  bekleyen_bagimlilik: number;
  bu_goreve_bagimli: number;
}

export interface GorevOlusturParams {
  baslik: string;
  aciklama?: string;
  oncelik?: GorevOncelik;
  proje_id?: string;
  parent_id?: string;
  son_tarih?: string;
  etiketler?: string;
}

export interface AltGorevOlusturParams {
  parent_id: string;
  baslik: string;
  aciklama?: string;
  oncelik?: GorevOncelik;
  son_tarih?: string;
  etiketler?: string;
}

export interface GorevUstDegistirParams {
  id: string;
  yeni_parent_id?: string;
}

export interface GorevGuncelleParams {
  id: string;
  durum: GorevDurum;
}

export interface GorevDuzenleParams {
  id: string;
  baslik?: string;
  aciklama?: string;
  oncelik?: GorevOncelik;
  proje_id?: string;
  son_tarih?: string;
}

export interface GorevListeleParams {
  durum?: GorevDurum;
  tum_projeler?: boolean;
  sirala?: 'son_tarih_asc' | 'son_tarih_desc';
  filtre?: 'acil' | 'gecmis';
  etiket?: string;
}

export interface GorevSilParams {
  id: string;
  onay: boolean;
}

export interface BagimlilikEkleParams {
  kaynak_id: string;
  hedef_id: string;
  baglanti_tipi: string;
}