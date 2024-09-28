'use client';
import {
    getEventById,
    getEventInterestById,
    migrateInterestsToParticipantInterests,
} from '$app/(public)/components/lib/serverAction';
import { Box, Button, Link, Paper, Typography } from '@mui/material';
import Image from 'next/image';
import NextLink from 'next/link';
import AwakeDragons from 'public/interessedragene/2024AwakeDragons1_1.png';
import HappyDragons from 'public/interessedragene/2024HappyDragons1_1.png';
import SleepyDragons from 'public/interessedragene/2024SleepyDragons1_1.png';
import VeryHappyDragons from 'public/interessedragene/2024VeryHappyDragons1_1.png';

type Props = {
    id: string;
};

const InterestPage = async ({ id }: Props) => {
    console.log('interest page: ', id);

    // const event = await getEventById(id);
    // const poolInterests = await getEventInterestById(id);
    // console.log('poolInterests: ', poolInterests);
    const handliMoveInterests = async () => {
        await migrateInterestsToParticipantInterests();
    };
    return (
        <Box>
            <Typography variant="h1">Ønskeliste: {id}</Typography>
            <Button variant="contained" color="primary" onClick={handliMoveInterests}>
                Migration
            </Button>
            <Paper>
                <Typography variant="h2">Veldig intreset</Typography>
                <Image src={VeryHappyDragons} alt="Veldig interessert" width={100} height={60} />
                <Box sx={{ padding: '1rem' }}>
                    <Typography variant="h3">Kai Norman</Typography>
                    <Typography variant="h3">Ola Norman</Typography>
                </Box>
            </Paper>
        </Box>
    );
};

export default InterestPage;
