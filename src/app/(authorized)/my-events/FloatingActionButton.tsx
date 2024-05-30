'use client';
import { faPlus } from '@fortawesome/free-solid-svg-icons/faPlus';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import Fab from '@mui/material/Fab';
import { useRouter } from 'next/navigation';
import { createMyEventDoc } from './actions';
import { useEffect } from 'react';

type Props = { newDocumentId: string };

const FloatingActionButton = ({ newDocumentId }: Props) => {
    const router = useRouter();
    useEffect(() => {
        router.prefetch(`/event/create/${newDocumentId}`);
    }, []);

    const handleClick = async () => {
        // await createMyEventDoc(newDocumentId);
        createMyEventDoc(newDocumentId);
        router.push(`/event/create/${newDocumentId}`);
    };
    return (
        <Fab
            color="primary"
            sx={{
                position: 'fixed',
                bottom: '1rem',
                right: '0.5rem',
            }}
            aria-label="edit"
            onClick={handleClick}
        >
            <FontAwesomeIcon icon={faPlus} size="3x" />
        </Fab>
    );
};

export default FloatingActionButton;
