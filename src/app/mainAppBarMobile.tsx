'use client';
import { AppBar, Box, Drawer, IconButton, Menu, MenuItem, styled, Toolbar } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import { ComponentProps, useState } from 'react';
import HomeIcon from '@mui/icons-material/Home';
import Link from 'next/link';
import LoginIcon from '@mui/icons-material/Login';
import LogoutIcon from '@mui/icons-material/Logout';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ArticleIcon from '@mui/icons-material/Article';
import EditNoteIcon from '@mui/icons-material/EditNote';

type Props = {};
const ITEM_HEIGHT = 48;

const MainAppBarMobile = ({}: Props) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    const handleClick: ComponentProps<'button'>['onClick'] = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };
    return (
        <>
            <AppBar position="fixed" color="primary" sx={{ top: 'auto', bottom: 0 }}>
                <Toolbar>
                    <IconButton href="/" aria-label="home" component={Link}>
                        <HomeIcon />
                    </IconButton>
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
                            {true ?
                                <>
                                    <MenuItem component={Link} href="/">
                                        <LogoutIcon />
                                        Logg ut
                                    </MenuItem>
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
                                </>
                            :   <MenuItem component={Link} href="/login">
                                    <LoginIcon />
                                    Logg inn
                                </MenuItem>
                            }
                        </Menu>
                    </Box>
                </Toolbar>
            </AppBar>
        </>
    );
};

export default MainAppBarMobile;
