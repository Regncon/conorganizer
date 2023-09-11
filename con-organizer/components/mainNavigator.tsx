"use client";

import * as React from 'react';
import Box from '@mui/material/Box';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import RestoreIcon from '@mui/icons-material/Restore';
import FavoriteIcon from '@mui/icons-material/Favorite';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import PasswordIcon from '@mui/icons-material/Password';
import { LocalActivity } from '@mui/icons-material';
import { Dialog } from '@mui/material';
import Login from './Login';
import ForgotPassword from './ForgotPassword';

const MainNavigator = () => {
  const [value, setValue] = React.useState(0);
  const [choice, setChoice] = React.useState("");

  return (
    <Box sx={{bottom:0, position:"fixed", width:"100%"}}>
      <BottomNavigation
        showLabels
        value={value}
        onChange={(event, newValue) => {
          setValue(newValue);
        }}
      >
        <BottomNavigationAction label="Logg inn" icon={<AccountCircleIcon />} onClick={()=>setChoice("login")} />
        <BottomNavigationAction label="Glemt passord" icon={<PasswordIcon />} onClick={()=>setChoice("newpassword")} />
        {/* <BottomNavigationAction label="Kj&oslash;p billett" icon={<LocalActivity />} /> */}
        <Dialog open={!!choice}>
          {choice === "login" ? <Login setChoice={setChoice} /> : <ForgotPassword /> }
        </Dialog>
      </BottomNavigation>
    </Box>
  );
}

export default MainNavigator;