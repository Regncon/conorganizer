'use client';
import { Button, Box, Slider, sliderClasses, Typography } from '@mui/material';
import Link from 'next/link';
import HelpIcon from '@mui/icons-material/Help';
import { useState } from 'react';
const marks = [
    { value: 1, label: '🥱 Ikke interessert' },
    { value: 2, label: '😑 Litt interessert' },
    { value: 3, label: '😊 Interessert' },
    { value: 4, label: '🤩 Veldig interessert' },
];
type Props = {};

const InterestSelector = ({}: Props) => {
    const [interest, setInterest] = useState<number>(0);
    const incrementInterest = () => {
        if (interest === 3) {
            setInterest(0);
        } else {
            setInterest(interest + 1);
        }
    };

    return (
        <>
            <Button
                variant="contained"
                color="primary"
                fullWidth
                sx={{
                    fontSize: '1.2rem',
                    textTransform: 'none',
                    marginBlock: '1rem',
                    minHeight: '62px',
                    maxWidth: 'var(--slider-interest-width)',
                }}
                onClick={incrementInterest}
            >
                {marks[interest].label}
            </Button>
            <Box sx={{ padding: '0.35rem', marginBottom: '0.1rem', maxWidth: 'var(--slider-interest-width)' }}>
                <Slider
                    onChange={(e) => {
                        const target = e.target as HTMLInputElement;
                        setInterest(Number(target.value));
                    }}
                    sx={{
                        color: 'primary.main',
                        [`.${sliderClasses.rail}`]: { backgroundColor: '#3d3b3b', height: '1rem' },
                        [`.${sliderClasses.track}`]: { height: '1rem' },
                        [`.${sliderClasses.mark}`]: {
                            borderRadius: '50%',
                            outlineColor: 'primary.main',
                            outlineWidth: '0.8rem',
                            outlineStyle: 'solid',
                            outlineOffset: '-1px',
                            opacity: '1',
                        },
                        [`.${sliderClasses.markActive}`]: { backgroundColor: 'primary.main' },
                        [`.${sliderClasses.thumb}:before`]: { boxShadow: 'unset' },
                        [`.${sliderClasses.valueLabelCircle}`]: { display: 'none' },
                    }}
                    marks
                    value={interest}
                    min={0}
                    max={3}
                />
            </Box>
            <Typography sx={{ marginBottom: '0.8rem', fontSize: '1rem' }}>
                Dra baren over for å melde din interesse!
            </Typography>
            <Box>
                <Box sx={{ marginBottom: '3rem' }}>
                    <Box
                        component={Link}
                        href="/hjelppaamelding"
                        sx={{
                            display: 'inline-flex',
                            gap: '0.4rem',
                            paddingLeft: '0.5rem',
                            color: 'primary.main',
                            maxWidth: 'var(--slider-interest-width)',
                        }}
                    >
                        <HelpIcon sx={{ scale: '1.5', placeSelf: 'center' }} />
                        <Typography component="p">Forvirret? Les mer om påmeldingsystemet</Typography>
                    </Box>
                </Box>
            </Box>
        </>
    );
};

export default InterestSelector;
