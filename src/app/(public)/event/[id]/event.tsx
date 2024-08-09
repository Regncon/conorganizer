'use client';
import {
    Box,
    Chip,
    Slider,
    TextField,
    Typography,
    sliderClasses,
    useTheme,
    type SxProps,
    type Theme,
} from '@mui/material';
import Image from 'next/image';
import NavigateBefore from '@mui/icons-material/NavigateBefore';
import IconButton from '@mui/material/IconButton';
import blekksprut2 from '$public/blekksprut2.jpg';
import HelpIcon from '@mui/icons-material/Help';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { faChevronLeft, faUserSecret } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faScroll } from '@fortawesome/free-solid-svg-icons/faScroll';
import { ConEvent } from '$lib/types';
import { db } from '$lib/firebase/firebase';
import { onSnapshot, doc, type Unsubscribe } from 'firebase/firestore';

export const dynamic = 'force-static';

const marks = [
    { value: 1, label: '🥱 Ikke interessert' },
    { value: 2, label: '😑 Litt interessert' },
    { value: 3, label: '😊 Interessert' },
    { value: 4, label: '🤩 Veldig interessert' },
];

type Props = {
    id?: string;
    eventData?: ConEvent;
    editable?: boolean;
};

const initialState: ConEvent = {
    gameMaster: '',
    id: '',
    shortDescription: '',
    description: '',
    system: '',
    title: '',
    email: '',
    name: '',
    phone: '',
    gameType: '',
    participants: 0,
    unwantedFridayEvening: false,
    unwantedSaturdayMorning: false,
    unwantedSaturdayEvening: false,
    unwantedSundayMorning: false,
    moduleCompetition: false,
    childFriendly: false,
    possiblyEnglish: false,
    adultsOnly: false,
    volunteersPossible: false,
    lessThanThreeHours: false,
    moreThanSixHours: false,
    beginnerFriendly: false,
    additionalComments: '',
    createdAt: '',
    createdBy: '',
    updateAt: '',
    updatedBy: '',
    subTitle: '',
    published: false,
    puljeFridayEvening: false,
    puljeSaturdayMorning: false,
    puljeSaturdayEvening: false,
    puljeSundayMorning: false,
};

const MainEvent = ({ id, eventData, editable = false }: Props) => {
    const [data, setData] = useState<ConEvent>(eventData ?? initialState);
    const [isEditingTitle, setIsEditingTitle] = useState<boolean>(false);
    const [isEditingGameMaster, setIsEditingGameMaster] = useState<boolean>(false);
    const [isEditingSystem, setIsEditingSystem] = useState<boolean>(false);
    const [isEditingDescription, setIsEditingDescription] = useState<boolean>(false);
    const [isEditingShortDescription, setIsEditingShortDescription] = useState<boolean>(false);
    const [interest, setInterest] = useState<number>(0);

    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (id !== undefined && eventData === undefined) {
            unsubscribeSnapshot = onSnapshot(doc(db, 'events', id), (snapshot) => {
                setData((snapshot.data() as ConEvent | undefined) ?? initialState);
            });
        }
        return () => {
            unsubscribeSnapshot?.();
        };
    }, [id, eventData]);

    const {
        palette: {
            background: { default: themeBackgroundColor },
        },
    } = useTheme();

    const paragraphStyle: SxProps<Theme> = { margin: '1rem 0' };
    // const strongStyle: SxProps<Theme> = { fontWeight: 700 };

    return (
        <Box>
            <Box
                sx={{
                    display: 'grid',
                    '& > *': { gridColumn: '1 / 2', gridRow: '1 / 2' },
                }}
            >
                <Box
                    component={Image}
                    src={blekksprut2}
                    alt="noe alt-tekst"
                    sx={{ width: '100%', height: 'auto', maxWidth: '100%', aspectRatio: '3.3 / 2' }}
                    placeholder="blur"
                    loading="lazy"
                />
                <Box sx={{ background: `linear-gradient(0deg, ${themeBackgroundColor}, transparent)` }}>
                    <Box
                        sx={{ display: 'grid', gridTemplateRows: '2rem 1fr', height: '100%', wordBreak: 'break-word' }}
                    >
                        <IconButton sx={{ placeSelf: 'start' }}>
                            <FontAwesomeIcon icon={faChevronLeft} />
                        </IconButton>
                        {isEditingTitle ?
                            <TextField
                                name="title"
                                value={data.title}
                                onChange={(e) => setData({ ...data, title: e.target.value })}
                                onBlur={() => setIsEditingTitle(false)}
                                autoFocus
                                variant="outlined"
                                fullWidth
                                sx={{ placeSelf: 'end center', paddingBottom: '2.5rem' }}
                            />
                            : <Typography
                                variant="h1"
                                align="center"
                                sx={{ placeSelf: 'end center', paddingBottom: '2.5rem' }}
                                onClick={() => editable && setIsEditingTitle(true)}
                            >
                                {data.title || 'Tittel'}
                            </Typography>
                        }
                    </Box>
                </Box>
            </Box>
            <Box sx={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
                <Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
                    <Box component={FontAwesomeIcon} icon={faUserSecret} size="2x" sx={{ color: 'primary.main' }} />
                    <Box>
                        <Typography component="span" sx={{ color: 'primary.main' }}>
                            {data.icons?.includes('rollespill') ? 'Gamemaster' : 'Arrangør'}
                        </Typography>
                        {isEditingGameMaster ?
                            <TextField
                                name="gameMaster"
                                value={data.gameMaster}
                                onChange={(e) => setData({ ...data, gameMaster: e.target.value })}
                                onBlur={() => setIsEditingGameMaster(false)}
                                autoFocus
                                variant="outlined"
                                fullWidth
                                sx={{ marginTop: '0.5rem' }}
                            />
                            : <Typography variant="h2" onClick={() => editable && setIsEditingGameMaster(true)}>
                                {data.gameMaster || 'Navn'}
                            </Typography>
                        }
                    </Box>
                </Box>
                <Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
                    <Box component={FontAwesomeIcon} icon={faScroll} size="2x" sx={{ color: 'primary.main' }} />
                    <Box>
                        <Typography component="span" sx={{ color: 'primary.main' }}>
                            System
                        </Typography>
                        {isEditingSystem ?
                            <TextField
                                name="system"
                                value={data.system}
                                onChange={(e) => setData({ ...data, system: e.target.value })}
                                onBlur={() => setIsEditingSystem(false)}
                                autoFocus
                                variant="outlined"
                                fullWidth
                                sx={{ marginTop: '0.5rem' }}
                            />
                            : <Typography variant="h2" onClick={() => editable && setIsEditingSystem(true)}>
                                {data.system || 'System'}
                            </Typography>
                        }
                    </Box>
                </Box>
            </Box>
            <Box
                sx={{ display: 'flex', gap: '.5em', overflowX: 'auto', marginBottom: '4rem', paddingBottom: '0.35rem' }}
            >
                {data.icons?.map((tag) => <Chip label={tag} key={tag} color="primary" icon={<NavigateBefore />} />)}
            </Box>
            <Box
                sx={{
                    backgroundColor: 'primary.main',
                    color: 'primary.contrastText',
                    minHeight: '62px',
                    textAlign: 'center',
                    display: 'grid',
                    placeContent: 'center',
                    borderRadius: '0.2rem',
                    marginBottom: '1rem',
                }}
            >
                <Typography sx={paragraphStyle} component="p">
                    {marks[interest].label}
                </Typography>
            </Box>
            <Box sx={{ padding: '0.35rem', marginBottom: '0.1rem' }}>
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
                    defaultValue={0}
                    min={0}
                    max={3}
                />
            </Box>

            <Typography sx={{ marginBottom: '0.8rem', textAlign: 'center' }}>
                Dra baren over for å melde din interesse!
            </Typography>

            <Box
                component={Link}
                href="#"
                sx={{
                    display: 'inline-flex',
                    gap: '0.4rem',
                    marginBottom: '3rem',
                    paddingLeft: '0.5rem',
                    color: 'primary.main',
                }}
            >
                <HelpIcon sx={{ scale: '1.5', placeSelf: 'center' }} />
                <Typography component="p">Forvirret? Les mer om påmeldingssystemet</Typography>
            </Box>

            {isEditingShortDescription ?
                <TextField
                    name="shortDescription"
                    value={data.shortDescription}
                    onChange={(e) => setData({ ...data, shortDescription: e.target.value })}
                    onBlur={() => setIsEditingShortDescription(false)}
                    autoFocus
                    variant="outlined"
                    fullWidth
                    sx={{ marginBottom: '1rem' }}
                />
                : <Typography
                    sx={{ ...paragraphStyle, marginBottom: '1rem', textAlign: 'center' }}
                    onClick={() => editable && setIsEditingShortDescription(true)}
                >
                    {data.shortDescription || 'Kort beskrivelse'}
                </Typography>
            }
            {isEditingDescription ?
                <TextField
                    name="description"
                    value={data.description}
                    onChange={(e) => setData({ ...data, description: e.target.value })}
                    onBlur={() => setIsEditingDescription(false)}
                    autoFocus
                    variant="outlined"
                    fullWidth
                    multiline
                    rows={4}
                    sx={{ marginBottom: '1rem' }}
                />
                : <Typography
                    sx={{ ...paragraphStyle, marginBottom: '1rem', textAlign: 'center' }}
                    onClick={() => editable && setIsEditingDescription(true)}
                >
                    {data.description || 'Lang beskrivelse'}
                </Typography>
            }
        </Box>
    );
};

export default MainEvent;
