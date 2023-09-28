'use client';

import { useEffect, useState } from 'react';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import AddIcon from '@mui/icons-material/Add';
import MenuIcon from '@mui/icons-material/Menu';
import PasswordIcon from '@mui/icons-material/Password';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import { Dialog, SpeedDial, SpeedDialAction } from '@mui/material';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import Box from '@mui/material/Box';
import { auth } from '@/lib/firebase';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
import { useAuth } from './AuthProvider';
import EditDialog from './EditDialog';
import ForgotPassword from './ForgotPassword';
import Login from './Login';
import Signup from './Signup';

const MainNavigator = () => {
    const [value, setValue] = useState(0);
    const [choice, setChoice] = useState('');
    const user = useAuth();
    const [openAdd, setOpenAdd] = useState(false);

    const { userSettings } = useUserSettings(user?.uid);
    const [showAddButton, setShowAddButton] = useState<boolean>(false);

    useEffect(() => {
        setShowAddButton(userSettings?.admin && user ? true : false);
    }, [user, userSettings]);

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
                        sx={
                            showAddButton
                                ? {
                                      display: 'block',
                                      '& .MuiSpeedDialAction-staticTooltipLabel': {
                                          width: '23ch',
                                      },
                                  }
                                : { display: 'none' }
                        }
                        key="AddEvent"
                        icon={<AddIcon />}
                        tooltipTitle="legg til arrangement"
                        tooltipOpen
                        onClick={() => setOpenAdd(true)}
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
                        label="Lag konto"
                        icon={<PersonAddIcon />}
                        onClick={() => setChoice('signup')}
                    />
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
                </BottomNavigation>
            )}
            <Dialog open={!!choice}>
                {choice === 'login' ? (
                    <Login setChoice={setChoice} />
                ) : choice === 'signup' ? (
                    <Signup setChoice={setChoice} />
                ) : (
                    <ForgotPassword setChoice={setChoice} />
                )}
            </Dialog>
            <EditDialog open={openAdd} handleClose={() => setOpenAdd(false)} />
        </Box>
    );
};

export default MainNavigator;
