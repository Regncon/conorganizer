'use client';
import { useRouter } from 'next/navigation';
import { createMyEventDoc } from './actions';
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
        <Box
            onClick={() => {
                console.log('clicked box 4 more box more click');
            }}
        >
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
                <CardActionArea sx={{}}>
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
                        <Box sx={{ display: 'flex', justifyContent:'center' }}>
                            <AddIcon sx={{ fontSize: '10rem' }} />
                        </Box>
                    </Box>
                </CardActionArea>
            </Card>
        </Box>
    );
};
export default AddEventCard;
