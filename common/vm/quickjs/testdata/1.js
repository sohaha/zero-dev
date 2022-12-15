// export default {}


const genRandom = (min, max) => (Math.random() * (max - min + 1) | 0) + min;



exports.default = {
    now: new Date(),
    rand: genRandom(1, 100),
}