'use client';
import { getAllParticipants } from '$app/(public)/components/lib/serverAction';
import { Participant } from '$lib/types';
import SearchIcon from '@mui/icons-material/Search';
import { FormControl, Input, InputAdornment, InputLabel, Paper } from '@mui/material';
import { useEffect, useState } from 'react';
import Fuse from 'fuse.js';

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
            setFilteredParticipants(participants); // Initially, all participants are shown
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
            setFilteredParticipants(participants); // Reset to all participants if the search is cleared
        }
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

            {/* Display the filtered participants */}
            <ul>
                {filteredParticipants.map((participant) => (
                    <li key={participant.id}>
                        {participant.firstName} - {participant.lastName}
                    </li>
                ))}
            </ul>
        </Paper>
    );
};

export default ManuallyAssignPlayer;
