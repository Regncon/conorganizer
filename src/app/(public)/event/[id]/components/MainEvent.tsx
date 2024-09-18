'use client';
import {
    Box,
    Button,
    Chip,
    CircularProgress,
    Paper,
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
import blekksprut2 from '$public/blekksprut2.jpg';
import HelpIcon from '@mui/icons-material/Help';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { faUserSecret } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faScroll } from '@fortawesome/free-solid-svg-icons/faScroll';
import { type PoolEvent } from '$lib/types';
import { db } from '$lib/firebase/firebase';
import { onSnapshot, doc, type Unsubscribe } from 'firebase/firestore';
import { MuiMarkdown } from 'mui-markdown';
import NavigatePreviousLink from './ui/NavigatePreviousLink';
import NavigateNextLink from './ui/NavigateNextLink';
import BackButton from '$app/(authorized)/components/BackButton';

export const dynamic = 'force-static';

const marks = [
    { value: 1, label: '🥱 Ikke interessert' },
    { value: 2, label: '😑 Litt interessert' },
    { value: 3, label: '😊 Interessert' },
    { value: 4, label: '🤩 Veldig interessert' },
];

type Props = {
    id?: string;
    eventData?: PoolEvent;
    editable?: boolean;
    editDescription?: (edit: boolean) => void;
    prevNavigationId?: string;
    nextNavigationId?: string;
    parent?: boolean;
};

const MainEvent = ({
    id,
    eventData,
    editable = false,
    parent = false,
    editDescription,
    prevNavigationId,
    nextNavigationId,
}: Props) => {
    const [data, setData] = useState<PoolEvent | undefined>(eventData);
    const [isEditingTitle, setIsEditingTitle] = useState<boolean>(false);
    const [isEditingGameMaster, setIsEditingGameMaster] = useState<boolean>(false);
    const [isEditingSystem, setIsEditingSystem] = useState<boolean>(false);
    const [isEditingShortDescription, setIsEditingShortDescription] = useState<boolean>(false);
    const [interest, setInterest] = useState<number>(0);

    useEffect(() => {
        let unsubscribeSnapshot: Unsubscribe | undefined;
        if (id !== undefined && eventData === undefined) {
            unsubscribeSnapshot = onSnapshot(doc(db, `${parent ? 'events' : 'pool-events'}`, id), (snapshot) => {
                setData(snapshot.data() as PoolEvent | undefined);
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
    const incrementInterest = () => {
        if (interest === 3) {
            setInterest(0);
        } else {
            setInterest(interest + 1);
        }
    };

    return data ?
        <Paper
            elevation={1}
            sx={{
                maxWidth: 'var(--max-width)',
                margin: 'auto',
                '&': {
                    '--max-width': '430px',
                    '--arrow-size': '2rem',
                },
            }}
        >
            <Box
                sx={{
                    display: 'grid',
                    '& > *': { gridColumn: '1 / 2', gridRow: '1 / 2' },
                    '& > img': { width: '100%', height: '100%' },
                }}
            >
                <Image
                    src={data?.smallImageURL ?? blekksprut2}
                    width={320}
                    height={273}
                    sizes="100vw"
                    alt="noe alt-tekst"
                    loading="lazy"
                />
                <Box
                    sx={{
                        background: `linear-gradient(0deg, ${themeBackgroundColor}, transparent)`,
                    }}
                >
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateRows: '2rem 1fr',
                            height: '100%',
                            wordBreak: 'break-word',
                        }}
                    >
                        <Box sx={{ placeSelf: 'start' }}>
                            <BackButton />
                        </Box>
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
                                sx={{ corsor: 'pointer', placeSelf: 'end center', paddingBottom: '2.5rem' }}
                                onClick={() => editable && setIsEditingTitle(true)}
                            >
                                {data.title || 'Tittel'}
                            </Typography>
                        }
                    </Box>
                </Box>
            </Box>
            <Box sx={{ padding: '1rem' }}>
                <Box sx={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
                    <Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
                        <Box
                            component={FontAwesomeIcon}
                            icon={faUserSecret}
                            size="2x"
                            sx={{ color: 'primary.main' }}
                        />
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
                                    sx={{ corsor: 'pointer', marginTop: '0.5rem' }}
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
                    sx={{
                        display: 'flex',
                        gap: '.5em',
                        overflowX: 'auto',
                        paddingBottom: '0.35rem',
                    }}
                >
                    {data.icons?.map((tag) => (
                        <Chip label={tag} key={tag} color="primary" icon={<NavigateBefore />} />
                    ))}
                </Box>
                <Button
                    variant="contained"
                    color="primary"
                    fullWidth
                    sx={{ fontSize: '1.2rem', textTransform: 'none', marginBottom: '1rem', minHeight: '62px' }}
                    onClick={incrementInterest}
                >
                    <Typography>Virker ikke enda: </Typography>
                    {marks[interest].label}
                </Button>
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
                        value={interest}
                        min={0}
                        max={3}
                    />
                </Box>
                <Typography sx={{ marginBottom: '0.8rem', textAlign: 'center' }}>
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
                            }}
                        >
                            <HelpIcon sx={{ scale: '1.5', placeSelf: 'center' }} />
                            <Typography component="p">Forvirret? Les mer om påmeldingsystemet</Typography>
                        </Box>
                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: !prevNavigationId || !nextNavigationId ? '1fr' : '1fr 1fr',
                                placeItems: 'space-between',
                                marginBlockStart: '1rem',
                                position: 'relative',
                                width: '100%',
                            }}
                        >
                            <NavigatePreviousLink previousNavigationId={prevNavigationId} />
                            <NavigateNextLink nextNavigationId={nextNavigationId} />
                        </Box>
                    </Box>
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
                        sx={{ corsor: 'pointer', ...paragraphStyle, marginBottom: '1rem', textAlign: 'center' }}
                        onClick={() => editable && setIsEditingShortDescription(true)}
                    >
                        {data.shortDescription || 'Kort beskrivelse'}
                    </Typography>
                }
                <Box
                    sx={{ cursor: 'pointer' }}
                    onClick={() => editable && editDescription && editDescription(true)}
                >
                    <MuiMarkdown>{data.description || '# Lang beskrivelse'}</MuiMarkdown>
                </Box>
            </Box>
        </Paper>
        : <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100dvh' }}>
            <Typography variant="h1" marginInlineEnd={2}>
                Henter arrangementet
            </Typography>
            <CircularProgress sx={{ fontSize: '18px' }} />
        </Box>;
};

export default MainEvent;
