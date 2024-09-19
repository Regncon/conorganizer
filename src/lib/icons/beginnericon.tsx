'use client';
import React from 'react';
import { SvgIcon } from '@mui/material';

type Props = {
    color?: "inherit" | "primary" | "secondary" | "action" | "disabled" | "error" | "info" | "success" | "warning";
    size?: 'small' | 'medium' | 'large';
  }

const BeginnerIcon = ({ color = 'primary', size = 'medium' }:Props) => {
  return (
    <SvgIcon
      style={{ fontSize: size === 'small' ? "2rem" : size === 'large' ? "4rem" : 24 }}
      color={'primary'}
    >
    <svg
   id="Layer_1"
   viewBox="0 0 38.652497 39.599998"
   version="1.1"
   width="38.652496"
   height="39.599998">
  <g
     id="Layer_9"
     transform="translate(-68.317903,-50.970001)">
    <path
       d="m 84.37,85.16 c 0.62,-2.58 0.47,-5.22 0.67,-7.83 0.08,-1.02 -0.22,-1.48 -1.27,-1.53 -3.75,-0.2 -6.81,-1.86 -9.62,-4.27 -2.65,-2.28 -4.13,-5.19 -5.27,-8.37 -0.37,-1.03 -0.38,-2.12 -0.51,-3.18 -0.34,-2.81 0.99,-4.41 3.82,-4.25 2.91,0.17 5.75,0.85 8.26,2.48 0.9,0.58 1.84,1.12 2.8,1.61 2.99,1.54 4.5,4.28 5.44,7.26 1.39,4.41 1.5,9.02 1.31,13.6 -0.09,2.01 0.54,3.94 0.35,5.94 -0.12,1.26 -0.22,2.54 -1.4,3.24 -1.64,0.97 -3.21,0.33 -3.92,-1.53 -0.38,-1.01 -0.65,-2.06 -0.65,-3.16 z"
       id="path563" />
    <path
       d="m 102.93,51.27 c 3.55,-0.29 4.3,2.17 3.97,4.8 -0.63,5 -3.27,8.72 -7,11.78 -1.46,1.2 -3.09,2.27 -4.96,2.79 -2.4,0.66 -3.3,0 -3.36,-2.46 -0.06,-2.48 -1.48,-4.31 -2.87,-6.16 -1.66,-2.21 -1.57,-2.73 0.35,-4.74 3.79,-3.98 8.52,-5.69 13.87,-6 z"
       id="path565" />
  </g>
</svg>

    </SvgIcon>
  );
};

export default BeginnerIcon;