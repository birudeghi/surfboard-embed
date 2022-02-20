import React from 'react';
import { ButtonGroup, SIZE, SHAPE, MODE } from 'baseui/button-group';
import { Button } from 'baseui/button';
import { useStyletron } from 'baseui';

import "./configuration.scss";

const Configuration = props => {
    const { boards, onSelect, disabled } = props;
    const [selected, setSelected] = React.useState();
    const [css] = useStyletron();

    const buttons = boards => (boards.map(board => (
            <Button
                className={css({
                    fontWeight: 400,
                    })}
            >
                {board.name}
            </Button>
        ))
    );
    
    const handleSelect = index => {
        onSelect(index);
        setSelected(index);
    }

    return (
        <div className="surfboard-configuration" id="surfboard-configuration">
            <ButtonGroup 
                size={SIZE.mini} 
                shape={SHAPE.pill}
                mode={MODE.radio}
                selected={selected}
                onClick={(event, index) => {
                    handleSelect(index);
                  }}
                disabled={disabled}
            >
                {buttons(boards)}
            </ButtonGroup>
        </div>
    )
}

export default Configuration;