'use client';
import { Box, Card, CardActions, FormControl, Input, InputAdornment, InputLabel, Typography } from '@mui/material';
import Link from 'next/link';
import SearchIcon from '@mui/icons-material/Search';
import { useEffect, useState } from 'react';
import Fuse from 'fuse.js';
import { Participant } from '$lib/types';
import AddIcon from '@mui/icons-material/Add';
import ParticipantCard from './ParticipantCard';

type Props = {
    participants: Participant[];
};

const ParticipantsList = ({ participants }: Props) => {
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [filteredParticipants, setFilteredParticipants] = useState<Participant[]>(participants);

    const fuseOptions = {
        keys: ['firstName', 'lastName', 'ticketCategory', 'orderId'],
        threshold: 0.5,
    };

    useEffect(() => {
        const fuse = new Fuse(participants, fuseOptions);

        const result = searchQuery ? fuse.search(searchQuery).map((res) => res.item) : participants;
        setFilteredParticipants(result);
    }, [searchQuery, participants]);

    return (
        <Box>
            <Card sx={{ maxWidth: '430px', paddingTop: '1rem' }}>
                <CardActions>
                    <FormControl sx={{ width: '100%' }} variant="standard">
                        <InputLabel htmlFor="input-with-icon-adornment">Søk etter bilett</InputLabel>
                        <Input
                            id="input-with-icon-adornment"
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            endAdornment={
                                <InputAdornment position="start">
                                    <SearchIcon />
                                </InputAdornment>
                            }
                        />
                    </FormControl>
                </CardActions>
            </Card>
            <Link href="/admin/dashboard/participants/add" prefetch style={{ textDecoration: 'none' }}>
                <Card
                    sx={{
                        maxWidth: '430px',
                        minWidth: '306px',
                        marginTop: '1rem',
                        marginBottom: '1rem',
                        height: '4rem',
                        width: '100%',
                        backgroundSize: 'cover',
                        borderRadius: '1.75rem',
                        opacity: 0.7,
                    }}
                >
                    <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                        <AddIcon sx={{ fontSize: '4rem' }} />
                    </Box>
                </Card>
            </Link>
            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                    gap: '1rem',
                }}
            >
                {filteredParticipants.map((participant) => (
                    <ParticipantCard key={participant.id} participant={participant} />
                ))}
            </Box>
        </Box>
    );
};
export default ParticipantsList;
