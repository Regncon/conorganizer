'use client';
import IconButton, { iconButtonClasses } from '@mui/material/IconButton';
import NoteAddIcon from '@mui/icons-material/NoteAdd';
import { useRouter } from 'next/navigation';
import { createMyEventDoc } from './actions';
import { touchRippleClasses } from '@mui/material';
import { useEffect } from 'react';
type Props = { newDocumentId: string };

const NewEventButton = ({ newDocumentId }: Props) => {
    const router = useRouter();
    useEffect(() => {
        router.prefetch(`/event/create/${newDocumentId}`);
    }, []);

    const handleClick = async () => {
        await createMyEventDoc(newDocumentId);
        router.push(`/event/create/${newDocumentId}`);
    };

    return (
        <IconButton
            sx={{
                border: '1px solid',
                borderColor: 'secondary.contrastText',
                position: 'absolute',
                zIndex: '22',
                right: '0.3125rem',
                top: '-0.375rem',
                [`.${touchRippleClasses.ripple}`]: {
                    color: 'black',
                },
                '&:hover, &:focus, &': {
                    backgroundColor: 'secondary.contrastText',
                },
            }}
            onClick={handleClick}
        >
            <NoteAddIcon />
        </IconButton>
    );
};

export default NewEventButton;
