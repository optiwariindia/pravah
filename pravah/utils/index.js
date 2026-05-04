const requiredEnv = [
    "PORT",
    "MONGODB"
]
Array.prototype.joinWithAnd = function (lang = "en", style = "long") {
    const formatterAnd = new Intl.ListFormat(lang, { style, type: 'conjunction' });
    return formatterAnd.format(this)
}
export async function checkEnvironmentVariables() {
    const { env } = process;
    let missingEnv = requiredEnv.filter(e => !(e in env));
    if(missingEnv.length===0)return;
    const error = new Error(`
        Environment Variable${missingEnv.length > 1 ? "s" : ""} ${missingEnv.joinWithAnd()} ${missingEnv.length > 1 ? "are" : "is"} missing.
        `)
    throw error
}