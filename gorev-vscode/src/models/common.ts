// Common types and enums for Gorev domain models

export enum GorevDurum {
  Beklemede = 'beklemede',
  DevamEdiyor = 'devam_ediyor',
  Tamamlandi = 'tamamlandi',
}

export enum GorevOncelik {
  Dusuk = 'dusuk',
  Orta = 'orta',
  Yuksek = 'yuksek',
}

export enum TemplateKategori {
  Genel = 'Genel',
  Teknik = 'Teknik',
  Ozellik = 'Özellik',
  Arastirma = 'Araştırma',
  Bug = 'Bug',
  Dokumantasyon = 'Dokümantasyon',
}

export interface Timestamp {
  olusturma_tarihi: string;
  guncelleme_tarihi: string;
}

export interface MCPError {
  code: number;
  message: string;
  data?: any;
}

export interface MCPToolResult {
  content: string;
  isError: boolean;
}

export type BaglantiTip = 'engelliyor' | 'iliskili' | 'depends_on';