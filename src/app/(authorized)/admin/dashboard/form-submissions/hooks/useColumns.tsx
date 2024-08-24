import { useMemo } from 'react';
import type { FormSubmissionColDef } from '../types';
import { useRouter } from 'next/navigation';
import { faArrowUpRightFromSquare } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { GridActionsCellItem } from '@mui/x-data-grid';
import type { Route } from 'next';

export const useColumns = () => {
    const router = useRouter();
    const columns: FormSubmissionColDef = useMemo(
        () => [
            // { field: 'id', headerName: 'ID', width: 255,  },
            {
                field: 'name',
                headerName: 'Navn',
                width: 140,
                editable: true,
            },
            {
                field: 'title',
                headerName: 'Tittel',
                width: 140,
                editable: true,
            },
            {
                field: 'subTitle',
                headerName: 'Kort beskrivelse',
                width: 240,
                editable: true,
            },
            {
                field: 'isSubmitted',
                headerName: 'Innsendt',
                type: 'boolean',
                width: 110,
                editable: true,
            },
            {
                field: 'isRead',
                headerName: 'Lest',
                type: 'boolean',
                width: 110,
                editable: true,
            },
            {
                field: 'isAccepted',
                headerName: 'Godkjent',
                type: 'boolean',
                width: 120,
                editable: true,
            },
            {
                field: 'actions',
                type: 'actions',
                headerName: 'Gå til førehandsvising',
                width: 150,
                getActions: (params) => [
                    <GridActionsCellItem
                        icon={<FontAwesomeIcon icon={faArrowUpRightFromSquare} />}
                        label="førehandsvising"
                        onClick={(e) => {
                            router.push(`/admin/dashboard/form-submissions/preview/${params.id}` as Route);
                        }}
                    />,
                ],
            },
        ],
        [router]
    );

    return columns;
};
