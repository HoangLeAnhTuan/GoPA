export type VocabularyLanguage = "JP" | "EN";
export interface Vocabulary { id:string; language:VocabularyLanguage; word:string; reading:string | null; meaning:string; example_sentence:string | null; current_box:number; next_review_at:string; }
export interface VocabularyInput { language:VocabularyLanguage; word:string; reading:string | null; meaning:string; exampleSentence:string | null; }
