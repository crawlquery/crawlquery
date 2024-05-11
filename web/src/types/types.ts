
export interface Page {
    id: string;
    url: string;
    title: string;
    description: string;
    content: string;
}

export interface Result {
    page_id: string;
    score: number;
    page: Page;
}
