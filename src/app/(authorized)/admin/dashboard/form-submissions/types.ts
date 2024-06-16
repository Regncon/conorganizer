import type { GridColDef } from '@mui/x-data-grid';

export type FormSubmission = {
    id: string;
    name: string;
    title: string;
    subTitle: string;
    isRead: boolean;
    isAccepted: boolean;
    documentPath: string;
};
export type FormSubmissionColDef = GridColDef<FormSubmission>[];
