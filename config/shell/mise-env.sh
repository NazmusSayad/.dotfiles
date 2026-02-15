while IFS= read -r line; do
  [[ -z "$line" || "$line" == \#* ]] && continue
  key="${line%%=*}"
  val="${line#*=}"
  export "$key=$val"
done < <(MSYS2_ARG_CONV_EXCL="*" mise env --dotenv)
