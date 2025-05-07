import React from "react";
import {Loader} from "@mantine/core";

interface CenteredLoaderProps {
  color?: string;
  size?: string | number;
}

export const CenteredLoader: React.FC<CenteredLoaderProps> = ({
  color = "teal",
  size = "xl"
}) => {
  return (
    <div className="flex justify-center items-center w-full">
      <Loader color={color} size={size}/>
    </div>
  );
};
