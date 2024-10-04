'use client';
import { getAllParticipants } from '$app/(public)/components/lib/serverAction';
import { Participant } from '$lib/types';
import SearchIcon from '@mui/icons-material/Search';
import { FormControl, Input, InputAdornment, InputLabel, Paper, Button, Box, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import Fuse from 'fuse.js';
import { assignPlayer } from './lib/actions';

type Props = {
    poolEventId: string;
};

const ManuallyAssignPlayer = ({ poolEventId }: Props) => {
    const [participants, setParticipants] = useState<Participant[]>([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [filteredParticipants, setFilteredParticipants] = useState<Participant[]>([]);

    // Initialize the participants
    useEffect(() => {
        getAllParticipants().then((participants) => {
            setParticipants(participants);
        });
    }, []);

    // Fuse.js options
    const fuse = new Fuse(participants, {
        keys: ['firstName', 'lastName'], // Fields to search
        threshold: 0.3, // Adjust for the sensitivity of the search
    });

    // Handle the search query change
    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        const query = event.target.value;
        setSearchQuery(query);

        if (query.trim()) {
            const results = fuse.search(query);
            setFilteredParticipants(results.map((result) => result.item)); // Get the actual participants
        } else {
            setFilteredParticipants([]); // Clear the filtered participants if the search is empty
        }
    };

    const handleAssignPlayer = (participant: Participant) => {
        const participantId = participant.id;
        if (!participantId) return;
        assignPlayer(participantId, poolEventId, true, false, undefined);
    };

    const handleAssignGameMaster = (participant: Participant) => {
        const participantId = participant.id;
        if (!participantId) return;
        assignPlayer(participantId, poolEventId, true, true, undefined);
    };

    return (
        <Paper sx={{ backgroundColor: 'rgba(0,0,0,0.1)', padding: '1rem' }}>
            <FormControl variant="standard">
                <InputLabel htmlFor="input-with-icon-adornment">Søk etter deltager</InputLabel>
                <Input
                    id="input-with-icon-adornment"
                    value={searchQuery}
                    onChange={handleSearch}
                    endAdornment={
                        <InputAdornment position="start">
                            <SearchIcon />
                        </InputAdornment>
                    }
                />
            </FormControl>

            {/* Display the filtered participants only if searchQuery is not empty */}
            {searchQuery && (
                <ul>
                    {filteredParticipants.slice(0, 10).map((participant, index) => (
                        <li
                            key={participant.id}
                            style={{
                                backgroundColor: index % 2 === 0 ? 'rgba(0, 0, 0, 0.05)' : 'rgba(0, 0, 0, 0.1)',
                                padding: '0.5rem',
                                marginBottom: '0.5rem',
                                borderRadius: '4px',
                            }}
                        >
                            <Box display="flex" alignItems="center" justifyContent="flex-start">
                                <Typography component={'span'} sx={{ minWidth: '25rem' }}>
                                    {participant.firstName} {participant.lastName}
                                </Typography>
                                <Box>
                                    <Button
                                        variant="contained"
                                        size="small"
                                        color="primary"
                                        onClick={() => handleAssignPlayer(participant)}
                                        sx={{ marginRight: '0.5rem' }}
                                    >
                                        Assign Player
                                    </Button>
                                    <Button
                                        variant="contained"
                                        size="small"
                                        color="secondary"
                                        onClick={() => handleAssignGameMaster(participant)}
                                    >
                                        Assign Game Master
                                    </Button>
                                </Box>
                            </Box>
                        </li>
                    ))}
                </ul>
            )}
        </Paper>
    );
};

export default ManuallyAssignPlayer;
