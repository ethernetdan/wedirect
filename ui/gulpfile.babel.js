// IMPORTS

import babelify from 'babelify'
import browserify from 'browserify'
import gulp from 'gulp'
import autoprefixer from 'gulp-autoprefixer'
import changed from 'gulp-changed'
import connect from 'gulp-connect'
import fileinclude from 'gulp-file-include'
import htmlmin from 'gulp-htmlmin'
import imagemin from 'gulp-imagemin'
import minify from 'gulp-minify-css'
import plumber from 'gulp-plumber'
import sass from 'gulp-sass'
import sourcemaps from 'gulp-sourcemaps'
import uglify from 'gulp-uglify'
import assign from 'lodash.assign'
import notifier from 'node-notifier'
import buffer from 'vinyl-buffer'
import source from 'vinyl-source-stream'
import watchify from 'watchify'

// ERROR HANDLER

const onError = function(error) {
  notifier.notify({
    'title': 'Error',
    'message': 'Compilation failure.'
  })

  console.log(error)
  this.emit('end')
}

// HTML

gulp.task('html', () => {
  return gulp.src('src/html/**/*.html')
    .pipe(plumber({ errorHandler: onError }))
    .pipe(fileinclude({ prefix: '@' }))
    .pipe(gulp.dest('../'))
    .pipe(connect.reload())
})

// JS

gulp.task('js', () => {
  return gulp.src('src/js/**/*.js')
    .pipe(plumber({ errorHandler: onError }))
    .pipe(gulp.dest('dist'))
    .pipe(connect.reload())
})

// SASS

gulp.task('sass', () => {
  return gulp.src('src/sass/style.scss')
    .pipe(plumber({ errorHandler: onError }))
    .pipe(sass())
    .pipe(autoprefixer({ browsers: [ 'last 2 versions', 'ie >= 9', 'Android >= 4.1' ] }))
    .pipe(minify())
    .pipe(gulp.dest('dist'))
    .pipe(connect.reload())
})

// WATCH

gulp.task('watch', () => {
  gulp.watch('src/html/**/*.html', ['html'])
  gulp.watch('src/sass/**/*.scss', ['sass'])
  gulp.watch('src/js/**/*.js', ['js'])
})

// TASKS

gulp.task('build', ['html', 'sass', 'js'])
gulp.task('default', ['build', 'watch'])
