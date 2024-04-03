import { useState } from "react";
import { v4 as uuidv4 } from "uuid";

const url = "ws://localhost:8080/connect";
const ws = new WebSocket(url);
let myid = "";

interface Cell {
  column: number;
  row: number;
}

const size = 10;

const strings = {
  phases: {
    setup: "Setup phase",
  },
  notification: {
    players: {
      first: "Waiting for the first player",
      second: "Waiting for the second player"
    }
  }
};

export default function App() {
  const [myField, setMyField] = useState<number[][]>();
  const [enemyField, setEnemyField] = useState<number[][]>();
  const [currentPhase, setCurrentPhase] = useState<string>(strings.phases.setup);

  function handleClick(cell: Cell) {
    ws.send(cell.column * 10 + cell.row + " ");
  }

  ws.addEventListener("message", (event) => {
    console.log("Recieved message: ", event.data)
  })

  return (
    <div>
      <h1 className="phase">{currentPhase}</h1>
      <div className="wrapper">
        <div>
          {myField ? myField.map((arr, indexTop) => {
            return (
              <div className="row">
                {arr.map((val, indexLow) => {
                  const cell: Cell = {
                    column: indexTop,
                    row: indexLow,
                  };
                  return (
                    <div
                      className="cell"
                      id={uuidv4()}
                      onClick={() => handleClick(cell)}
                    >
                      {val}
                    </div>
                  );
                })}
              </div>
            );
          }) : <p className="notification">{strings.notification.players.first}</p>}
        </div>
        <div>
          {enemyField ? enemyField.map((arr, indexTop) => {
            return (
              <div className="row">
                {arr.map((val, indexLow) => {
                  const cell: Cell = {
                    column: indexTop,
                    row: indexLow,
                  };
                  return (
                    <div
                      className="cell"
                      id={uuidv4()}
                      onClick={() => handleClick(cell)}
                    >
                      {val}
                    </div>
                  );
                })}
              </div>
            );
          }) : <p className="notification">{strings.notification.players.second}</p>}
        </div>
      </div>
    </div>
  );
}

function oneDimensionToTwoDimensions(arr: number[], size: number) {
  let arr1: number[] = [];
  const arr2: number[][] = [];
  for (const i in arr) {
    arr1.push(arr[i]);
    if (+i % size === size - 1) {
      arr2.push(arr1);
      arr1 = [];
    }
  }
  return arr2;
}
