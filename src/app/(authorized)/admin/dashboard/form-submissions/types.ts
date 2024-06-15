import type { GridColDef } from '@mui/x-data-grid';

export type FormSubmission = {
    id: string;
    title: string;
    subTitle: string;
    isRead: boolean;
    isAccepted: boolean;
};
export type FormSubmissionColDef = GridColDef<FormSubmission>[];
