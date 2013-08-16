package chunker


func calculateChunks(bucket, dir string, r Range) []*Chunk {
  if r.To == 0 {
    return []*Chunk{&Chunk{From: r.From,}}
  }

  chunks := make([]*Chunk, count)
  for i := r.From; i <= r.To; i += CHUNK_SIZE {
    n := int(math.Floor(float64(i) / float64(CHUNK_SIZE)))
    from := i
    if i != r.From { from = CHUNK_SIZE * n }

    to := (n + 1) * CHUNK_SIZE - 1
    if to > r.To { to = r.To }
    key := bucket + "_" + strconv.Itoa(n)

    chunk := &Chunk{
      From: from,
      To: to,
      Key: key,
      DateFile: dir + key,
    }
    chunks = append(chunks, chunk)
  }
  return chunks
}
