'use strict';

const gulp = require('gulp');
const autoprefixer = require('gulp-autoprefixer')
const connect = require('gulp-connect')
const plumber = require('gulp-plumber')
const sass = require('gulp-sass')
const notifier = require('node-notifier')

const onError = function(error) {
  notifier.notify({
    'title': 'Error',
    'message': 'Compilation failure.'
  })

  console.log(error)
  this.emit('end')
}

gulp.task('html', () => {
  return gulp.src('src/html/**/*.html')
    .pipe(plumber({ errorHandler: onError }))
    .pipe(gulp.dest('dist'))
    .pipe(connect.reload())
})

gulp.task('js', () => {
  return gulp.src('src/js/**/*.js')
    .pipe(plumber({ errorHandler: onError }))
    .pipe(gulp.dest('dist'))
    .pipe(connect.reload())
})

gulp.task('sass', () => {
  return gulp.src('src/sass/style.scss')
    .pipe(plumber({ errorHandler: onError }))
    .pipe(sass())
    .pipe(autoprefixer({ browsers: [ 'last 2 versions', 'ie >= 9', 'Android >= 4.1' ] }))
    .pipe(gulp.dest('dist'))
    .pipe(connect.reload())
})

gulp.task('watch', () => {
  gulp.watch('src/html/**/*.html', ['html'])
  gulp.watch('src/sass/**/*.scss', ['sass'])
  gulp.watch('src/js/**/*.js', ['js'])
})

gulp.task('build', ['html', 'sass', 'js'])
gulp.task('default', ['build', 'watch'])
