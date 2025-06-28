import { GorevDurum, GorevOncelik, Timestamp } from './common';

export { GorevDurum, GorevOncelik } from './common';

export interface Gorev extends Timestamp {
  id: string;
  baslik: string;
  aciklama: string;
  durum: GorevDurum;
  oncelik: GorevOncelik;
  proje_id: string;
  son_tarih?: string;
  etiketler?: string[];
  bagimliliklar?: Bagimlilik[];
}

export interface GorevDetay extends Gorev {
  proje_isim: string;
  bagimliliklar?: Bagimlilik[];
}

export interface Bagimlilik {
  kaynak_id: string;
  hedef_id: string;
  baglanti_tip: string;
  hedef_baslik?: string;
  hedef_durum?: GorevDurum;
}

export interface GorevOlusturParams {
  baslik: string;
  aciklama?: string;
  oncelik?: GorevOncelik;
  proje_id?: string;
  son_tarih?: string;
  etiketler?: string;
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