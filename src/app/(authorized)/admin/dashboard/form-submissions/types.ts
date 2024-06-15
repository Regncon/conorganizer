import type { GridColDef } from '@mui/x-data-grid';

export type FromSubmission = {
    id: string;
    title: string;
    subtitle: string;
    isRead: boolean | string;
    isAccepted: boolean | string;
};
export type FromSubmissionColDef = GridColDef<FromSubmission>[];
