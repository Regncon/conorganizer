import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { Card, CardContent, CardHeader } from '@mui/material';
import Image from 'next/image';
import rook from '$lib/image/rook.svg';
import type { EventCardProps } from './types';

export default function EventCardBig({ title, gameMaster, system, icons }: Omit<EventCardProps, 'shortDescription'>) {
    return (
        <Card
            sx={{
                backgroundImage: 'url(/blekksprut2.jpg)',
                maxHeight: '187px',
                maxWidth: '149.5px',
                height: '100%',
                width: '100%',
                backgroundSize: 'cover',
                borderRadius: '1.75rem',
            }}
        >
            <CardHeader sx={{ height: '79px', alignItems: 'flex-end', padding: '1rem' }} />
            <CardContent
                sx={{
                    height: '126px',
                    backgroundColor: 'rgba(0,0,0,0.5)',
                    backdropFilter: 'blur(4px)',
                    padding: '0rem',
                    display: 'grid',
                    placeContent: 'center',
                }}
            >
                <Typography sx={{ fontSize: '15px' }}> {gameMaster} </Typography>
                <Typography sx={{ fontSize: '17px', fontWeight: 'bold', paddingBottom: '0.25rem' }}>{title}</Typography>
                <Typography
                    sx={{
                        color: 'secondary.contrastText',
                        paddingBottom: '0.15rem',
                    }}
                >
                    {system}
                </Typography>
                <Box sx={{ display: 'flex', gap: '1rem' }}>
                    <Box component={Image} priority src={rook} alt="rook icon" />
                    <Box component={Image} priority src={rook} alt="rook icon" />
                    <Box component={Image} priority src={rook} alt="rook icon" />
                    <Box component={Image} priority src={rook} alt="rook icon" />
                </Box>
            </CardContent>
        </Card>
    );
}
