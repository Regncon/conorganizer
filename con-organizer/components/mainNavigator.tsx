'use client';

import { useContext, useState } from 'react';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import MenuIcon from '@mui/icons-material/Menu';
import PasswordIcon from '@mui/icons-material/Password';
import { Dialog, SpeedDial, SpeedDialAction } from '@mui/material';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import Box from '@mui/material/Box';
import { auth } from '@/lib/firebase';
import { AuthContext, useAuth } from './AuthProvider';
import ForgotPassword from './ForgotPassword';
import Login from './login';

const MainNavigator = () => {
    const [value, setValue] = useState(0);
    const [choice, setChoice] = useState('');
    const user = useAuth();

    function logout() {
        return auth.signOut();
    }

    return (
        <Box sx={{ bottom: 0, position: 'fixed', width: '100%' }}>
            {user ? (
                <SpeedDial
                    ariaLabel="Alternativer"
                    sx={{ position: 'absolute', bottom: 16, right: 16 }}
                    icon={<MenuIcon />}
                >
                    <SpeedDialAction
                        key="logout"
                        icon={<AccountCircleIcon />}
                        tooltipTitle="logg ut"
                        sx={{
                            '& .MuiSpeedDialAction-staticTooltipLabel': {
                                width: '10ch',
                            },
                        }}
                        onClick={logout}
                        tooltipOpen
                    />
                    <SpeedDialAction
                        sx={{
                            '& .MuiSpeedDialAction-staticTooltipLabel': {
                                width: '16ch',
                            },
                        }}
                        key="changePwd"
                        icon={<PasswordIcon />}
                        tooltipTitle="endre passord"
                        onClick={() => setChoice('newpassword')}
                        tooltipOpen
                    />
                </SpeedDial>
            ) : (
                <BottomNavigation
                    showLabels
                    value={value}
                    onChange={(event, newValue) => {
                        setValue(newValue);
                    }}
                >
                    <BottomNavigationAction
                        label="Logg inn"
                        icon={<AccountCircleIcon />}
                        onClick={() => setChoice('login')}
                    />
                    <BottomNavigationAction
                        label="Glemt passord"
                        icon={<PasswordIcon />}
                        onClick={() => setChoice('newpassword')}
                    />
                    {/* <BottomNavigationAction label="Kj&oslash;p billett" icon={<LocalActivity />} /> */}
                    <Dialog open={!!choice}>
                        {choice === 'login' ? (
                            <Login setChoice={setChoice} />
                        ) : (
                            <ForgotPassword setChoice={setChoice} />
                        )}
                    </Dialog>
                </BottomNavigation>
            )}
        </Box>
    );
};

export default MainNavigator;
