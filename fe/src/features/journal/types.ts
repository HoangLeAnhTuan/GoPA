export interface Journal { id:string; title:string; content:string; tags:string[]; created_at:string; updated_at:string; }
export interface JournalInput { title:string; content:string; tags:string[]; }
