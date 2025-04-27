/* eslint-disable @typescript-eslint/no-explicit-any */
'use client';

import React, {useState} from 'react';

import { GridCol } from '@mantine/core';
import { Input, InputBase, Combobox, useCombobox } from '@mantine/core';

export default function Search({
  options, onChange
}: { onChange: any, options: any[]}) {
  const [value, setValue] = useState<string | null>(null);
  const combobox = useCombobox({
    onDropdownClose: () => combobox.resetSelectedOption(),
  });

  const comboboxOptions = options.map((item) => (
    <Combobox.Option value={item.name} key={item.name}>
      {item.name}
     </Combobox.Option>
  ));
  return (
    <GridCol span={4}>
      <Combobox
        store={combobox}
        onOptionSubmit={(val) => {
          setValue(val);
          onChange(val)
          combobox.closeDropdown();
        }}>
      <Combobox.Target>
          <InputBase
            component="button"
            type="button"
            pointer
            rightSection={<Combobox.Chevron />}
            rightSectionPointerEvents="none"
            onClick={() => combobox.toggleDropdown()}
          >
            {value || <Input.Placeholder>Pick an option</Input.Placeholder>}
          </InputBase>
        </Combobox.Target>
        <Combobox.Dropdown>
          <Combobox.Options>{comboboxOptions}</Combobox.Options>
        </Combobox.Dropdown>
      </Combobox>
    </GridCol>
  );
}
