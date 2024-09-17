import { useMemo } from 'react';
import type { FormSubmissionColDef } from '../lib/types';
import { useRouter } from 'next/navigation';
import { faArrowUpRightFromSquare } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import type { Route } from 'next';
import { Link, Tooltip } from '@mui/material';

export const useColumns = () => {
    const router = useRouter();
    const columns: FormSubmissionColDef = useMemo(
        () => [
            // { field: 'id', headerName: 'ID', width: 255,  },
            {
                field: 'name',
                headerName: 'Navn',
                width: 140,
            },
            {
                field: 'title',
                headerName: 'Tittel',
                width: 140,
            },
            {
                field: 'subTitle',
                headerName: 'Kort beskrivelse',
                width: 240,
            },
            {
                field: 'isSubmitted',
                headerName: 'Innsendt',
                type: 'boolean',
                width: 110,
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
            },
            {
                field: 'actions',
                type: 'actions',
                headerName: 'Gå til førehandsvising',
                width: 150,
                getActions: (params) => {
                    return [
                        <Tooltip title="førehandsvising">
                            <Link
                                sx={{ color: 'white' }}
                                href={
                                    `/admin/dashboard/form-submissions/preview/${params.id}/${params.row.userId}` as Route
                                }
                            >
                                <FontAwesomeIcon icon={faArrowUpRightFromSquare} />
                            </Link>
                        </Tooltip>,
                    ];
                },
            },
        ],
        [router]
    );

    return columns;
};
