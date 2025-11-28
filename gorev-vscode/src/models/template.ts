import { TemplateKategori } from './common';

export { TemplateKategori } from './common';

export interface GorevTemplate {
  id: string;
  isim: string;
  tanim: string;
  varsayilan_baslik: string;
  aciklama_template: string;
  alanlar: TemplateAlan[];
  ornek_degerler: Record<string, string>;
  kategori: TemplateKategori;
  aktif: boolean;
}

export interface TemplateAlan {
  isim: string;
  tur: 'metin' | 'sayi' | 'tarih' | 'secim';
  zorunlu: boolean;
  varsayilan?: string;
  secenekler?: string[];
}

export interface TemplateOlusturParams {
  template_id: string;
  degerler: Record<string, unknown>;
}

export interface TemplateListeleParams {
  kategori?: TemplateKategori;
}