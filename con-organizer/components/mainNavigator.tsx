'use client';

import { useContext, useState } from 'react';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import PasswordIcon from '@mui/icons-material/Password';
import { Dialog } from '@mui/material';
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
    console.log('useAuth', useAuth());
    console.log(user);

    function logout() {
        return auth.signOut();
    }

    return (
        <Box sx={{ bottom: 0, position: 'fixed', width: '100%' }}>
            {user ? (
                <BottomNavigation
                    showLabels
                    value={value}
                    onChange={(event, newValue) => {
                        setValue(newValue);
                    }}
                >
                    <BottomNavigationAction label="Logg ut" icon={<AccountCircleIcon />} onClick={logout} />
                    <BottomNavigationAction
                        label="Endre passord"
                        icon={<PasswordIcon />}
                        onClick={() => setChoice('newpassword')}
                    />
                    {/* <BottomNavigationAction label="Kj&oslash;p billett" icon={<LocalActivity />} /> */}
                    <Dialog open={!!choice}>
                        <ForgotPassword setChoice={setChoice} />
                    </Dialog>
                </BottomNavigation>
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
