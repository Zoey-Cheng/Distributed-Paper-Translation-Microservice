import SparkMD5 from "spark-md5";

const arrayBufferToString = (exportedPrivateKey: ArrayBuffer) => {
  const byteArray = new Uint8Array(exportedPrivateKey);
  let byteString = '';
  for (var i = 0; i < byteArray.byteLength; i++) {
    byteString += String.fromCodePoint(byteArray[i]);
  }
  return byteString;
}

const getFileHash = async (file: Blob) => {
  return new Promise((resolve: (md5: string) => void, reject) => {
    const reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = function () {
      const buf = reader.result as ArrayBuffer
      resolve(SparkMD5.hashBinary(arrayBufferToString(buf)))
    }
    reader.onerror = (ev) => {
      reject("error")
    }
  })
}

export default getFileHash;