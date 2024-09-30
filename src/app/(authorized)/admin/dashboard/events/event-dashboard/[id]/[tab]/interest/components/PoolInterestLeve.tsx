// 'use client';
import { poolTitles } from '$app/(authorized)/admin/dashboard/rooms/components/lib/helpers';
import {
    getEventById,
    getEventInterestById,
    migrateInterestsToParticipantInterests,
} from '$app/(public)/components/lib/serverAction';
import { InterestLevel } from '$lib/enums';
import { Interest } from '$lib/types';
import { Box, Button, Link, Paper, Typography } from '@mui/material';
import Image from 'next/image';
import { interestLevelToImage, InterestLevelToLabel } from './lib/helpers/InterestHelper';

type Props = {
    interestLevel: InterestLevel;
    poolInterest: Interest[];
};

const PoolInterestLevel = ({ interestLevel, poolInterest }: Props) => {
    return (
        <Paper elevation={2} sx={{ paddingBlock: '0.5rem' }}>
            <Box sx={{ display: 'flex', gap: '1rem' }}>
                <Image
                    src={interestLevelToImage[interestLevel]}
                    alt={InterestLevelToLabel[interestLevel]}
                    width={100}
                    height={50}
                />
                <Typography variant="h3">{InterestLevelToLabel[interestLevel]}</Typography>
            </Box>

            {poolInterest
                .filter((interest) => interest.interestLevel === interestLevel)
                .map((veryInterested) => (
                    <Box component="ul" key={veryInterested.id}>
                        <Typography component="li">
                            {veryInterested.participantFirstName} {veryInterested.participantLastName}
                        </Typography>
                    </Box>
                ))}
        </Paper>
    );
};
export default PoolInterestLevel;
