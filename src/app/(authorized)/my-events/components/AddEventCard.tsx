'use client';
import { useRouter } from 'next/navigation';
import { createMyEventDoc } from '../lib/actions';
import { useEffect } from 'react';
import AddIcon from '@mui/icons-material/Add';
import { Box, Card, CardActionArea, CardHeader, Paper, Typography } from '@mui/material';

type Props = { newDocumentId: string };

const AddEventCard = ({ newDocumentId }: Props) => {
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
        <Card
            sx={{
                minHeight: '267px',
                minWidth: '306px',
                maxHeight: '267px',
                maxWidth: '306px',
                height: '100%',
                width: '100%',
                backgroundSize: 'cover',
                borderRadius: '1.75rem',
                opacity: 0.7,
            }}
        >
            <CardActionArea onClick={handleClick}>
                <Box
                    sx={{
                        display: 'block',
                        minHeight: '267px',
                        minWidth: '306px',
                        maxHeight: '267px',
                        maxWidth: '306px',
                        height: '100%',
                        width: '100%',
                    }}
                >
                    <Typography variant="h3" sx={{ textAlign: 'center' }}>
                        Legg til nytt arrangement
                    </Typography>
                    <Box sx={{ display: 'flex', justifyContent: 'center', paddingBlockStart: '1.5rem' }}>
                        <AddIcon sx={{ fontSize: '10rem' }} />
                    </Box>
                </Box>
            </CardActionArea>
        </Card>
    );
};
export default AddEventCard;
