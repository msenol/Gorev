import { Timestamp } from './common';

export interface Proje extends Timestamp {
  id: string;
  isim: string;
  tanim: string;
  gorev_sayisi?: number;
  tamamlanan_sayisi?: number;
  devam_eden_sayisi?: number;
  bekleyen_sayisi?: number;
}

export interface ProjeOlusturParams {
  isim: string;
  tanim?: string;
}

export interface AktifProje {
  proje_id: string;
  proje?: Proje;
}