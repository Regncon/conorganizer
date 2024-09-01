'use client';
import { AppBar, Box, Button, IconButton, Menu, MenuItem, Toolbar, type SxProps, type Theme } from '@mui/material';
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
import { useState, type ComponentProps } from 'react';
type Props = {
    sx?: SxProps<Theme>;
    user: boolean;
    admin: boolean;
    mobile: boolean;
};

const MainAppBarContent = ({ sx, user, mobile, admin }: Props) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);

    const handleClick: ComponentProps<'button'>['onClick'] = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const desktopBottomContent = (
        <>
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
        </>
    );
    const mobileBottomContent = (
        <>
            <Box sx={{ flexGrow: 1 }} />
            <IconButton href="/" aria-label="home" component={Link}>
                <Box sx={{ flexGrow: 1 }} />
                <HomeIcon />
            </IconButton>

            <IconButton href="/?filters=favorites" aria-label="favorites" component={Link}>
                <FavoriteIcon />
            </IconButton>
            <IconButton href="/?filters=assigned" aria-label="assigned-events" component={Link}>
                <GroupsIcon />
            </IconButton>
        </>
    );

    return (
        <AppBar position="fixed" color="primary" sx={sx}>
            <Toolbar sx={{ gap: 2 }}>
                {mobile ? mobileBottomContent : desktopBottomContent}
                <Box>
                    {user ?
                        <>
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
};

export default MainAppBarContent;
