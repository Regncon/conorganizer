'use client';
import { AppBar, Box, Button, IconButton, Menu, MenuItem, Toolbar, Typography } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import { ComponentProps, useEffect, useState } from 'react';
import HomeIcon from '@mui/icons-material/Home';
import Link from 'next/link';
import LoginIcon from '@mui/icons-material/Login';
import LogoutIcon from '@mui/icons-material/Logout';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ArticleIcon from '@mui/icons-material/Article';
import EditNoteIcon from '@mui/icons-material/EditNote';
import { onAuthStateChanged, type User } from 'firebase/auth';
import { firebaseAuth } from '$lib/firebase/firebase';

type Props = {};

const MainAppBarDesktop = ({ }: Props) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);

    const [user, setUser] = useState<User | null>();
    useEffect(() => {
        const unsubscribeUser = onAuthStateChanged(firebaseAuth, (user) => {
            setUser(user);
        });

        return () => {
            unsubscribeUser();
        };
    }, [user]);

    const handleClick: ComponentProps<'button'>['onClick'] = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };
    return (
        <>
            <AppBar position="fixed" color="primary" sx={{ display: { xs: 'none', sm: 'block' } }}>
                <Toolbar>
                    <Button startIcon={<HomeIcon />} href="/" aria-label="home" component={Link}>
                        <Box sx={{ flexGrow: 1 }} />
                        Hjem
                    </Button>
                    <Box sx={{ flexGrow: 1 }} />
                    <Box>
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
                        {user ?
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
                                <MenuItem component={Link} href="/my-events">
                                    <ArticleIcon />
                                    Mine arrangementer
                                </MenuItem>
                                <MenuItem component={Link} href="/">
                                    <AccountCircleIcon />
                                    Min profil
                                </MenuItem>
                                <MenuItem component={Link} href="/admin">
                                    <AdminPanelSettingsIcon />
                                    Admin
                                </MenuItem>
                                <MenuItem component={Link} href="/admin/dashboard/events">
                                    <EditNoteIcon />
                                    Rediger arrangementer
                                </MenuItem>
                                <MenuItem component={Link} href="/logout">
                                    <LogoutIcon />
                                    Logg ut
                                </MenuItem>
                            </Menu>
                            : <Link href="/login">
                                <LoginIcon />
                                <Typography component={'span'}>Logg inn</Typography>
                            </Link>
                        }
                    </Box>
                </Toolbar>
            </AppBar>
        </>
    );
};

export default MainAppBarDesktop;
