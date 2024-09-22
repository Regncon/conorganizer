'use client';
import {
    AppBar,
    Box,
    Button,
    IconButton,
    Menu,
    MenuItem,
    Toolbar,
    Typography,
    type SxProps,
    type Theme,
} from '@mui/material';
import HomeIcon from '@mui/icons-material/Home';
import MenuIcon from '@mui/icons-material/Menu';
import LoginIcon from '@mui/icons-material/Login';
import LogoutIcon from '@mui/icons-material/Logout';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ArticleIcon from '@mui/icons-material/Article';
import EditNoteIcon from '@mui/icons-material/EditNote';
import FavoriteIcon from '@mui/icons-material/Favorite';
import GroupsIcon from '@mui/icons-material/Groups';
import Link from 'next/link';
import { forwardRef, useEffect, useState, type ComponentProps } from 'react';
import ParticipantSelector from '$ui/participant/ParticipantSelector';
import { collection, Firestore, getDocs } from 'firebase/firestore';
import { Participant } from '$lib/types';
import { db } from '$lib/firebase/firebase';

type Props = {
    sx?: SxProps<Theme>;
    isLoggedIn: boolean;
    admin: boolean;
    mobile: boolean;
};

const GetAllParticipants = async (db: Firestore) => {
    if (db === null) {
        const response = {
            type: 'error',
            message: 'Ikke autorisert',
            error: 'getAuthorizedAuth failed',
        };
        throw response;
    }

    const querySnapshot = await getDocs(collection(db, 'participants'));
    querySnapshot.forEach((doc) => {
        console.log(doc.id, ' => ', doc.data());
    });
    return querySnapshot.docs.map((doc) => doc.data()) as Participant[];
};

const MainAppBarContent = forwardRef<HTMLElement, Props>(({ sx, isLoggedIn, admin, mobile }, ref) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);

    const handleClick: ComponentProps<'button'>['onClick'] = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const [participants, setParticipants] = useState<Participant[] | null>(null);
    useEffect(() => {
        GetAllParticipants(db).then((participants) => {
            setParticipants(participants);
        });
    }, []);

    const desktopBottomContent = (
        <>
            <Typography variant="h1" sx={{}}>
                BETA TEST
            </Typography>
            <Button startIcon={<HomeIcon />} href="/" aria-label="home" component={Link}>
                <Box sx={{ flexGrow: 1 }} />
                Hjem
            </Button>
            <Button
                sx={{ color: 'white' }}
                startIcon={<FavoriteIcon />}
                href="/?filters=favorites"
                aria-label="favorites"
                component={Link}
            >
                Favoritter
            </Button>
            <Button
                sx={{ color: 'white' }}
                startIcon={<GroupsIcon />}
                href="/?filters=assigned"
                aria-label="assigned-events"
                component={Link}
            >
                Påmeldt
            </Button>
            <Box sx={{ flexGrow: 1 }} />
            {isLoggedIn ?
                <>
                    {participants ?
                        <ParticipantSelector participants={participants} />
                    :   <Button href="/my-profile/my-tickets">Hent billett</Button>}

                    <IconButton
                        aria-label="more"
                        id="long-button"
                        aria-controls={open ? 'long-menu' : undefined}
                        aria-expanded={open ? 'true' : undefined}
                        aria-haspopup="true"
                        onClick={handleClick}
                    >
                        <MenuIcon />
                    </IconButton>
                </>
            :   null}
        </>
    );
    const mobileBottomContent = (
        <>
            <Box sx={{ flexGrow: 1 }} />
            <IconButton href="/" aria-label="home" component={Link}>
                <HomeIcon fontSize="large" />
            </IconButton>
            <IconButton href="/?filters=favorites" aria-label="favorites" component={Link}>
                <FavoriteIcon fontSize="large" />
            </IconButton>
            <IconButton href="/?filters=assigned" aria-label="assigned-events" component={Link}>
                <GroupsIcon fontSize="large" />
            </IconButton>
            {isLoggedIn ?
                <IconButton
                    aria-label="more"
                    id="long-button"
                    aria-controls={open ? 'long-menu' : undefined}
                    aria-expanded={open ? 'true' : undefined}
                    aria-haspopup="true"
                    onClick={handleClick}
                >
                    <MenuIcon fontSize="large" />
                </IconButton>
            :   null}
        </>
    );

    return (
        <AppBar position="fixed" color="primary" sx={sx} ref={ref}>
            <Toolbar sx={{ gap: 2, placeItems: 'center', height: '100%' }}>
                {mobile ? mobileBottomContent : desktopBottomContent}
                <Box>
                    {isLoggedIn ?
                        <>
                            <Menu
                                id="long-menu"
                                MenuListProps={{
                                    'aria-labelledby': 'long-button',
                                }}
                                anchorEl={anchorEl}
                                open={open}
                                onClose={handleClose}
                                onClick={handleClose}
                            >
                                <MenuItem key="logout" component={Link} href="/logout">
                                    <LogoutIcon />
                                    Logg ut
                                </MenuItem>
                                <MenuItem key="my-events" component={Link} href="/my-events">
                                    <ArticleIcon />
                                    Mine arrangementer
                                </MenuItem>
                                <MenuItem key="my-profile" component={Link} href="/my-profile">
                                    <AccountCircleIcon />
                                    Min profil
                                </MenuItem>
                                {admin ?
                                    [
                                        <MenuItem key="admin" component={Link} href="/admin">
                                            <AdminPanelSettingsIcon />
                                            Admin
                                        </MenuItem>,
                                        <MenuItem key="admin-dashboard" component={Link} href="/admin/dashboard/events">
                                            <EditNoteIcon />
                                            Rediger arrangementer
                                        </MenuItem>,
                                    ]
                                :   null}
                            </Menu>
                        </>
                    :   <Button component={Link} href="/login">
                            <LoginIcon />
                            Logg inn
                        </Button>
                    }
                </Box>
                {mobile ?
                    <Box sx={{ flexGrow: 1 }} />
                :   null}
            </Toolbar>
        </AppBar>
    );
});

export default MainAppBarContent;
